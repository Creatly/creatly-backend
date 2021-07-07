package service

import (
	"context"
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

// TODO implement background workers for removing bad/broken files

const (
	_workersCount   = 3
	_workerInterval = time.Minute
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
	// select for processing
	file, err := s.repo.GetForUploading(ctx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}

		return err
	}

	logger.Infof("processing file %s", file.Name)

	// read file
	f, err := os.Open(file.Name)
	if err != nil {
		return err
	}

	// upload
	url, err := s.storage.Upload(ctx, storage.UploadInput{
		File:        f,
		Size:        file.Size,
		ContentType: file.ContentType,
		Name:        s.generateFilenameFromFile(file),
	})
	if err != nil {
		return err
	}

	if err := s.repo.UpdateStatusAndSetURL(ctx, file.ID, url); err != nil {
		return err
	}

	return nil
}

func (s *FilesService) Upload(ctx context.Context, inp UploadInput) (string, error) {
	return s.storage.Upload(ctx, storage.UploadInput{
		File:        inp.File,
		Size:        inp.Size,
		ContentType: inp.ContentType,
		Name:        s.generateFilename(inp),
	})
}

func (s *FilesService) GetByID(ctx context.Context, id, schoolId primitive.ObjectID) (domain.File, error) {
	return s.repo.GetByID(ctx, id, schoolId)
}

// TODO leave only one method.
func (s *FilesService) generateFilename(inp UploadInput) string {
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), getFileExtension(inp.Filename))
	folder := folders[inp.Type]

	return fmt.Sprintf("%s/%s/%s/%s", s.env, inp.SchoolID.Hex(), folder, filename)
}

func (s *FilesService) generateFilenameFromFile(file domain.File) string {
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), getFileExtension(file.Name))
	folder := folders[file.Type]

	fileNameParts := strings.Split(file.Name, "-") // first part is schoolId

	return fmt.Sprintf("%s/%s/%s/%s", s.env, fileNameParts[0], folder, filename)
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")

	return parts[len(parts)-1]
}
