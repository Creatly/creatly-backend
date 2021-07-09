package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
)

type adminGetVideoResponse struct {
	Status domain.FileStatus `json:"status"`
	URL    string            `json:"url"`
}

// @Summary Get Video By ID
// @Security AdminAuth
// @Tags admins-media
// @Description  get video by id
// @ModuleID adminGetVideo
// @Accept  json
// @Produce  json
// @Param id path string true "video id"
// @Success 200 {object} adminGetVideoResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/media/videos/{id} [get]
func (h *Handler) adminGetVideo(c *gin.Context) {
	id, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	file, err := h.services.Files.GetByID(c.Request.Context(), id, school.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, adminGetVideoResponse{
		Status: file.Status,
		URL:    file.URL,
	})
}
