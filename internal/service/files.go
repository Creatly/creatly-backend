package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"github.com/zhashkevych/creatly-backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
	"github.com/zhashkevych/creatly-backend/pkg/storage"
)

const (
	_workersCount   = 2
	_workerInterval = time.Second * 10
)

var folders = map[domain.FileType]string{
	domain.Image: "images",
	domain.Video: "videos",
	domain.Other: "other",
}

type FilesService struct {
	repo    repository.Files
	storage storage.Provider
	env     string
}

func NewFilesService(repo repository.Files, storage storage.Provider, env string) *FilesService {
	return &FilesService{repo: repo, storage: storage, env: env}
}

func (s *FilesService) Save(ctx context.Context, file domain.File) (primitive.ObjectID, error) {
	return s.repo.Create(ctx, file)
}

func (s *FilesService) UpdateStatus(ctx context.Context, fileName string, status domain.FileStatus) error {
	return s.repo.UpdateStatus(ctx, fileName, status)
}

func (s *FilesService) GetByID(ctx context.Context, id, schoolId primitive.ObjectID) (domain.File, error) {
	return s.repo.GetByID(ctx, id, schoolId)
}

func (s *FilesService) UploadAndSaveFile(ctx context.Context, file domain.File) (string, error) {
	defer removeFile(file.Name)

	file.UploadStartedAt = time.Now()

	if _, err := s.Save(ctx, file); err != nil {
		return "", err
	}

	return s.upload(ctx, file)
}

func (s *FilesService) InitStorageUploaderWorkers(ctx context.Context) {
	for i := 0; i < _workersCount; i++ {
		go s.processUploadToStorage(ctx)
	}
}

func (s *FilesService) processUploadToStorage(ctx context.Context) {
	for {
		if err := s.uploadToStorage(ctx); err != nil {
			logger.Error("uploadToStorage(): ", err)
		}

		time.Sleep(_workerInterval)
	}
}

func (s *FilesService) uploadToStorage(ctx context.Context) error {
	file, err := s.repo.GetForUploading(ctx)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}

		return err
	}

	defer removeFile(file.Name)

	logger.Infof("processing file %s", file.Name)

	url, err := s.upload(ctx, file)
	if err != nil {
		if err := s.repo.UpdateStatus(ctx, file.Name, domain.StorageUploadError); err != nil {
			return err
		}

		return err
	}

	logger.Infof("file %s processed successfully", file.Name)

	if err := s.repo.UpdateStatusAndSetURL(ctx, file.ID, url); err != nil {
		return err
	}

	return nil
}

func (s *FilesService) upload(ctx context.Context, file domain.File) (string, error) {
	f, err := os.Open(file.Name)
	if err != nil {
		return "", err
	}

	info, _ := f.Stat()
	logger.Infof("file info: %+v", info)

	defer f.Close()

	return s.storage.Upload(ctx, storage.UploadInput{
		File:        f,
		Size:        file.Size,
		ContentType: file.ContentType,
		Name:        s.generateFilename(file),
	})
}

func (s *FilesService) generateFilename(file domain.File) string {
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), getFileExtension(file.Name))
	folder := folders[file.Type]

	fileNameParts := strings.Split(file.Name, "-") // first part is schoolId

	return fmt.Sprintf("%s/%s/%s/%s", s.env, fileNameParts[0], folder, filename)
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")

	return parts[len(parts)-1]
}

func removeFile(filename string) {
	if err := os.Remove(filename); err != nil {
		logger.Error("removeFile(): ", err)
	}
}
