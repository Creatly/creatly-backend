package v1

import (
	"net/http"

	"github.com/zhashkevych/creatly-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

type schoolSettingsResponse struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Settings    domain.Settings `json:"settings"`
}

// @Summary School GetSettings
// @Tags school-settings
// @Description school get settings
// @ModuleID getSchoolSettings
// @Produce  json
// @Success 200 {object} domain.Settings
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /settings [get]
func (h *Handler) getSchoolSettings(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, schoolSettingsResponse{
		Name:        school.Name,
		Description: school.Description,
		Settings:    school.Settings,
	})
}
