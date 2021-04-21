package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/zhashkevych/courses-backend/pkg/storage"
)

type FilesService struct {
	storage storage.Provider
	env     string
}

func NewFilesService(storage storage.Provider, env string) *FilesService {
	return &FilesService{storage: storage, env: env}
}

func (s *FilesService) Upload(ctx context.Context, inp UploadInput) (string, error) {
	filename := uuid.New().String()
	filepath := fmt.Sprintf("%s/%s/%s", s.env, inp.SchoolID.Hex(), filename)

	return s.storage.Upload(ctx, storage.UploadInput{
		File:        inp.File,
		Size:        inp.Size,
		ContentType: inp.ContentType,
		Name:        filepath,
	})
}
