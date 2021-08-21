package v1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zhashkevych/creatly-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
)

const (
	maxUploadSize = 5 << 20 // 5 megabytes
	maxVideoSize  = 2 << 30 // 2 gigabytes
)

type contentRange struct {
	rangeStart int64
	rangeEnd   int64
	fileSize   int64
}

func (cr *contentRange) parse(c *gin.Context) error {
	contentRangeHeader := c.Request.Header.Get("Content-Range")
	rangeAndSizeNumbers := strings.Split(contentRangeHeader, " ")
	rangeAndSize := strings.Split(rangeAndSizeNumbers[1], "/")
	rangeParts := strings.Split(rangeAndSize[0], "-")

	var err error

	cr.rangeStart, err = strconv.ParseInt(rangeParts[0], 10, 64)
	if err != nil {
		return err
	}

	cr.rangeEnd, err = strconv.ParseInt(rangeParts[1], 10, 64)
	if err != nil {
		return err
	}

	cr.fileSize, err = strconv.ParseInt(rangeAndSize[1], 10, 64)
	if err != nil {
		return err
	}

	return nil
}

func (cr *contentRange) isUploadCompleted() bool {
	return cr.fileSize == cr.rangeEnd
}

func (cr *contentRange) initialUploadRequest() bool {
	return cr.rangeStart == 0
}

func (cr *contentRange) chunkSize() int64 {
	return cr.rangeEnd - cr.rangeStart
}

func (cr *contentRange) validSize(maxSize int64) bool {
	return cr.fileSize <= maxSize
}

var (
	imageTypes = map[string]interface{}{
		"image/jpeg": nil,
		"image/png":  nil,
	}

	videoTypes = map[string]interface{}{
		"video/mp4":                 nil,
		"application/octet-stream":  nil,
		"text/plain; charset=utf-8": nil, // for strange files with such content-type
	}
)

type uploadResponse struct {
	URL string `json:"url"`
}

// @Summary Admin upload Image
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
func (h *Handler) adminUploadImage(c *gin.Context) { //nolint:funlen
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

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

	tempFilename := fmt.Sprintf("%s-%s", school.ID.Hex(), fileHeader.Filename)

	f, err := os.OpenFile(tempFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to create temp file")

		return
	}

	defer f.Close()

	if _, err := io.Copy(f, bytes.NewReader(buffer)); err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to write chunk to temp file")

		return
	}

	url, err := h.services.Files.UploadAndSaveFile(c.Request.Context(), domain.File{
		SchoolID:    school.ID,
		Type:        domain.Image,
		ContentType: contentType,
		Name:        tempFilename,
		Size:        fileHeader.Size,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, &uploadResponse{url})
}

type uploadVideoResponse struct {
	ID string `json:"id"`
}

// @Summary Admin upload Video
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
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	rangeInfo := new(contentRange)
	if err := rangeInfo.parse(c); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if !rangeInfo.validSize(maxVideoSize) {
		newResponse(c, http.StatusBadRequest, "file is bigger than 2 gigabytes")

		return
	}

	logger.Infof("chunk size: %d", rangeInfo.chunkSize())

	buffer := make([]byte, rangeInfo.chunkSize())

	if _, err = file.Read(buffer); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	contentType := http.DetectContentType(buffer)
	logger.Infof("chunk content type: %s", contentType)

	if _, ex := videoTypes[contentType]; !ex {
		newResponse(c, http.StatusBadRequest, "file type is not supported")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	// todo strip symbols in filename
	tempFilename := fmt.Sprintf("%s-%s", school.ID.Hex(), fileHeader.Filename)

	f, err := os.OpenFile(tempFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to create temp file")

		return
	}

	defer f.Close()

	if _, err := io.Copy(f, bytes.NewReader(buffer)); err != nil {
		if err := h.services.Files.UpdateStatus(c.Request.Context(), tempFilename, domain.ClientUploadError); err != nil {
			logger.Errorf("failed to write chunk to temp file & failed to update file status: %s", err.Error())
		}

		if err := os.Remove(tempFilename); err != nil {
			logger.Errorf("failed to delete corrupted temp file: %s", err.Error())
		}

		newResponse(c, http.StatusInternalServerError, "failed to write chunk to temp file")

		return
	}

	if rangeInfo.initialUploadRequest() {
		id, err := h.services.Files.Save(c.Request.Context(), domain.File{
			Type:            domain.Video,
			Status:          domain.ClientUploadInProgress,
			Name:            tempFilename,
			Size:            rangeInfo.fileSize,
			UploadStartedAt: time.Now(),
			ContentType:     contentType,
			SchoolID:        school.ID,
		})
		if err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to save file info to DB")

			return
		}

		c.JSON(http.StatusCreated, &uploadVideoResponse{ID: id.Hex()})

		return
	}

	if rangeInfo.isUploadCompleted() {
		if err := h.services.Files.UpdateStatus(c.Request.Context(), tempFilename, domain.UploadedByClient); err != nil {
			newResponse(c, http.StatusInternalServerError, "failed to update file status")

			return
		}
	}

	c.Status(http.StatusOK)
}

// @Summary Admin upload File
// @Security AdminAuth
// @Tags admins-upload
// @Description admin upload file
// @ModuleID adminUploadFile
// @Accept mpfd
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} uploadResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/upload/file [post]
func (h *Handler) adminUploadFile(c *gin.Context) { //nolint:funlen
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

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

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	tempFilename := fmt.Sprintf("%s-%s", school.ID.Hex(), fileHeader.Filename)

	f, err := os.OpenFile(tempFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to create temp file")

		return
	}

	defer f.Close()

	if _, err := io.Copy(f, bytes.NewReader(buffer)); err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to write chunk to temp file")

		return
	}

	url, err := h.services.Files.UploadAndSaveFile(c.Request.Context(), domain.File{
		SchoolID:    school.ID,
		Type:        domain.Other,
		ContentType: contentType,
		Name:        tempFilename,
		Size:        fileHeader.Size,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, &uploadResponse{url})
}
