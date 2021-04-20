package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/pkg/storage"
	"math/rand"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fileNameLength = 16
)

type FilesService struct {
	storage storage.Provider
}

func NewFilesService(storage storage.Provider) *FilesService {
	return &FilesService{storage: storage}
}

func (s *FilesService) Upload(ctx context.Context, inp UploadInput) (string, error) {
	return s.storage.Upload(ctx, storage.UploadInput{
		File:        inp.File,
		Size:        inp.Size,
		ContentType: inp.ContentType,
		Name:        generateFileName(),
	})
}

func generateFileName() string {
	b := make([]byte, fileNameLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
