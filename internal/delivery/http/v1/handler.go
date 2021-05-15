package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/service"
	"github.com/zhashkevych/creatly-backend/pkg/auth"
)

type Handler struct {
	services          *service.Services
	tokenManager      auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:          services,
		tokenManager:      tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initCoursesRoutes(v1)
		h.initStudentsRoutes(v1)
		h.initCallbackRoutes(v1)
		h.initAdminRoutes(v1)

		v1.GET("/settings", h.setSchoolFromRequest, h.getSchoolSettings)
		v1.GET("/promocodes/:code", h.setSchoolFromRequest, h.getPromo)
	}
}
