package service

import (
	"context"
	"fmt"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"github.com/google/uuid"
	"github.com/zhashkevych/creatly-backend/pkg/storage"
)

const (
	uploadTimeout = time.Minute
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

func (s *FilesService) Upload(ctx context.Context, inp UploadInput) (string, error) {
	ctx, clFn := context.WithTimeout(ctx, uploadTimeout)
	defer clFn()

	return s.storage.Upload(ctx, storage.UploadInput{
		File:        inp.File,
		Size:        inp.Size,
		ContentType: inp.ContentType,
		Name:        s.generateFilename(inp),
	})
}

func (s *FilesService) generateFilename(inp UploadInput) string {
	filename := fmt.Sprintf("%s.%s", uuid.New().String(), inp.FileExtension)
	folder := folders[inp.Type]

	return fmt.Sprintf("%s/%s/%s/%s", s.env, inp.SchoolID.Hex(), folder, filename)
}
