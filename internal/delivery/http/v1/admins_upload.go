package v1

import (
	"bytes"
	"net/http"
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
		"video/mp4":                nil,
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
func (h *Handler) adminUploadVideo(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxVideoUploadSize)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	defer file.Close()

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

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
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

func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	return parts[len(parts)-1]
}
