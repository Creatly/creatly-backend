package v1

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/service"
	mock_service "github.com/zhashkevych/creatly-backend/internal/service/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler_adminUploadImage(t *testing.T) {
	type mockBehavior func(r *mock_service.MockFiles, filepath, extension, contentType string, fileSize int64) error

	school := domain.School{
		ID: primitive.NewObjectID(),
	}

	tests := []struct {
		name         string
		filePath     string
		contentType  string
		extension    string
		fileSize     int64
		returnUrl    string
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:        "Ok jpg",
			filePath:    "./fixtures/image.jpg",
			fileSize:    434918,
			contentType: "image/jpeg",
			extension:   "jpg",
			mockBehavior: func(r *mock_service.MockFiles, filepath, extension, contentType string, fileSize int64) error {
				file, err := os.Open(filepath)
				if err != nil {
					return err
				}

				defer file.Close()

				buffer := make([]byte, fileSize)
				_, err = file.Read(buffer)
				if err != nil {
					return err
				}

				r.EXPECT().UploadAndSaveFile(context.Background(), domain.File{
					Type:        domain.Image,
					Name:        fmt.Sprintf("%s-image.jpg", school.ID.Hex()),
					ContentType: contentType,
					Size:        fileSize,
					SchoolID:    school.ID,
				}).Return("https://storage/image.jpg", nil)

				return nil
			},
			statusCode:   200,
			responseBody: `{"url":"https://storage/image.jpg"}`,
		},
		{
			name:        "Ok png",
			filePath:    "./fixtures/image.png",
			fileSize:    33764,
			contentType: "image/png",
			extension:   "png",
			mockBehavior: func(r *mock_service.MockFiles, filepath, extension, contentType string, fileSize int64) error {
				file, err := os.Open(filepath)
				if err != nil {
					return err
				}

				defer file.Close()

				buffer := make([]byte, fileSize)
				_, err = file.Read(buffer)
				if err != nil {
					return err
				}

				r.EXPECT().UploadAndSaveFile(context.Background(), domain.File{
					Type:        domain.Image,
					Name:        fmt.Sprintf("%s-image.png", school.ID.Hex()),
					ContentType: contentType,
					Size:        fileSize,
					SchoolID:    school.ID,
				}).Return("https://storage/image.png", nil)

				return nil
			},
			statusCode:   200,
			responseBody: `{"url":"https://storage/image.png"}`,
		},
		{
			name:     "Image too large",
			filePath: "./fixtures/large.jpeg",
			mockBehavior: func(r *mock_service.MockFiles, filepath, extension, contentType string, fileSize int64) error {
				return nil
			},
			statusCode:   400,
			responseBody: `{"message":"http: request body too large"}`,
		},
		{
			name:     "PDF upload",
			filePath: "./fixtures/ccc.pdf",
			mockBehavior: func(r *mock_service.MockFiles, filepath, extension, contentType string, fileSize int64) error {
				return nil
			},
			statusCode:   400,
			responseBody: `{"message":"file type is not supported"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			s := mock_service.NewMockFiles(c)
			err := tt.mockBehavior(s, tt.filePath, tt.extension, tt.contentType, tt.fileSize)
			assert.NoError(t, err)

			services := &service.Services{Files: s}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/upload", func(c *gin.Context) {
				c.Set(schoolCtx, school)
			}, handler.adminUploadImage)

			// Create Request
			file, err := os.Open(tt.filePath)
			assert.NoError(t, err)

			defer file.Close()

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))

			_, err = io.Copy(part, file)
			assert.NoError(t, err)

			err = writer.Close()
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/upload", body)
			req.Header.Add("Content-Type", writer.FormDataContentType())

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)

			// Remove files
			filenameParts := strings.Split(tt.filePath, "/")
			filename := fmt.Sprintf("%s-%s", school.ID.Hex(), filenameParts[len(filenameParts)-1])
			os.Remove(filename)
		})
	}
}
