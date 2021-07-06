package v1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/service"
)

const (
	maxImageUploadSize = 5 << 20 // 5 megabytes
	maxVideoUploadSize = 2 << 30 // 2 gigabytes
)

var (
	imageTypes = map[string]interface{}{
		"image/jpeg": nil,
		"image/png":  nil,
	}

	videoTypes = map[string]interface{}{
		"application/octet-stream": nil,
	}
)

type uploadResponse struct {
	URL string `json:"url"`
}

// @Summary Admin Upload Image
// @Security AdminAuth
// @Tags admins-upload
// @Description admin upload image
// @ModuleID adminUploadImage
// @Accept mpfd
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} uploadResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/upload/image [post]
func (h *Handler) adminUploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxImageUploadSize)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	defer file.Close()

	buffer := make([]byte, fileHeader.Size)

	if _, err := file.Read(buffer); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	contentType := http.DetectContentType(buffer)

	// Validate File Type
	if _, ex := imageTypes[contentType]; !ex {
		newResponse(c, http.StatusBadRequest, "file type is not supported")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	url, err := h.services.Files.Upload(c.Request.Context(), service.UploadInput{
		Type:          service.FileTypeImage,
		File:          bytes.NewBuffer(buffer),
		FileExtension: getFileExtension(fileHeader.Filename),
		ContentType:   contentType,
		Size:          fileHeader.Size,
		SchoolID:      school.ID,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, &uploadResponse{url})
}

// @Summary Admin Upload Video
// @Security AdminAuth
// @Tags admins-upload
// @Description admin upload video
// @ModuleID adminUploadVideo
// @Accept mpfd
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} uploadResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/upload/video [post]
func (h *Handler) adminUploadVideo(c *gin.Context) { //nolint:funlen
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxVideoUploadSize)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	tempFilename := fmt.Sprintf("%s-%s", school.ID, fileHeader.Filename)

	f, err := os.OpenFile(tempFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModeAppend)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to create temp file")

		return
	}

	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to write chunk to temp file")

		return
	}

	completed, err := isFileUploadCompleted(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	// TODO implement background uploading to Spaces ?
	// 1. Create upload ID and save it to Mongo
	// 2. Run background worker that will manage uploads and save URL to DB
	// 3. Implement endpoint to retrieve URL by ID
	if completed {
		buffer := make([]byte, fileHeader.Size)

		if _, err = file.Read(buffer); err != nil {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		contentType := http.DetectContentType(buffer)
		if _, ex := videoTypes[contentType]; !ex {
			newResponse(c, http.StatusBadRequest, "file type is not supported")

			return
		}

		url, err := h.services.Files.Upload(c.Request.Context(), service.UploadInput{
			Type:          service.FileTypeVideo,
			File:          bytes.NewBuffer(buffer),
			FileExtension: getFileExtension(fileHeader.Filename),
			ContentType:   "video/mp4",
			Size:          fileHeader.Size,
			SchoolID:      school.ID,
		})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())

			return
		}

		c.JSON(http.StatusOK, &uploadResponse{url})
	}

	c.Status(http.StatusOK)
}

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")

	return parts[len(parts)-1]
}

func isFileUploadCompleted(c *gin.Context) (bool, error) {
	contentRangeHeader := c.Request.Header.Get("Content-Range")
	rangeAndSize := strings.Split(contentRangeHeader, "/")
	rangeParts := strings.Split(rangeAndSize[0], "-")

	rangeMax, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return false, err
	}

	fileSize, err := strconv.Atoi(rangeAndSize[1])
	if err != nil {
		return false, err
	}

	return fileSize == rangeMax, nil
}
