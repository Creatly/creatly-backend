package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/service"
	"github.com/zhashkevych/creatly-backend/pkg/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initCoursesRoutes(v1)
		h.initStudentsRoutes(v1)
		h.initCallbackRoutes(v1)
		h.initAdminRoutes(v1)

		v1.GET("/settings", h.setSchoolFromRequest, h.getSchoolSettings)
		v1.GET("/promocodes/:code", h.setSchoolFromRequest, h.getPromo)
		v1.GET("/offers/:id", h.setSchoolFromRequest, h.getOffer)
	}
}

func parseIdFromPath(c *gin.Context, param string) (primitive.ObjectID, error) {
	idParam := c.Param(param)
	if idParam == "" {
		return primitive.ObjectID{}, errors.New("empty id param")
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid id param")
	}

	return id, nil
}
