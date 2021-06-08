package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/zhashkevych/creatly-backend/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileType int

const (
	FileTypeImage FileType = iota
	FileTypeVideo
)

const (
	uploadTimeout = time.Minute
)

var folders = map[FileType]string{
	FileTypeImage: "images",
	FileTypeVideo: "videos",
}

type (
	UploadInput struct {
		File          io.Reader
		FileExtension string
		Size          int64
		ContentType   string
		SchoolID      primitive.ObjectID
		Type          FileType
	}

	FilesService struct {
		storage storage.Provider
		env     string
	}
)

func NewFilesService(storage storage.Provider, env string) *FilesService {
	return &FilesService{storage: storage, env: env}
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
