package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go"
	"strings"
	"time"
)

type StorageType int

const (
	StorageSpaces StorageType = iota
	StorageS3
)

const (
	timeout = time.Second * 5
)

type FileStorage struct {
	client      *minio.Client
	bucket      string
	endpoint    string
	storageType StorageType
}

func NewFileStorage(client *minio.Client, bucket, endpoint string, storageType StorageType) *FileStorage {
	return &FileStorage{
		client:      client,
		bucket:      bucket,
		endpoint:    endpoint,
		storageType: storageType,
	}
}

// todo: image compression
func (fs *FileStorage) Upload(ctx context.Context, input UploadInput) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType:  input.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	ctx, clFn := context.WithTimeout(ctx, timeout)
	defer clFn()

	_, err := fs.client.PutObjectWithContext(ctx,
		fs.bucket, input.Name, input.File, input.Size, opts)
	if err != nil {
		return "", err
	}

	return fs.generateFileURL(input.Name), nil
}

func (fs *FileStorage) generateFileURL(fileName string) string {
	if fs.storageType == StorageSpaces {
		return generateSpacesLink(fs.bucket, fs.endpoint, fileName)
	}

	return generateS3Link(fs.bucket, fs.endpoint, fileName)
}

func generateSpacesLink(bucket, endpoint, filename string) string {
	return fmt.Sprintf("https://%s.%s/%s", bucket, endpoint, filename)
}

func generateS3Link(bucket, endpoint, filename string) string {
	endpoint = strings.Replace(endpoint, "localstack", "localhost", -1)
	return fmt.Sprintf("http://%s/%s/%s", endpoint, bucket, filename)
}
