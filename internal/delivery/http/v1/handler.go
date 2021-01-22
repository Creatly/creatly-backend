package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/service"
)

type Handler struct {
	schoolsService  service.Schools
	studentsService service.Students
}

func NewHandler(schoolsService service.Schools, studentsService service.Students) *Handler {
	return &Handler{schoolsService: schoolsService, studentsService: studentsService}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initStudentsRoutes(v1)
	}
}
