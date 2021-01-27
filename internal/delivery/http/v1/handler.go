package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/service"
	"github.com/zhashkevych/courses-backend/pkg/auth"
)

type Handler struct {
	schoolsService  service.Schools
	studentsService service.Students
	coursesService  service.Courses
	ordersService   service.Orders
	paymentsService service.Payments
	tokenManager    auth.TokenManager
}

func NewHandler(schoolsService service.Schools, studentsService service.Students, coursesService service.Courses, ordersService service.Orders,
	paymentsService service.Payments, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		schoolsService:  schoolsService,
		studentsService: studentsService,
		coursesService:  coursesService,
		ordersService:   ordersService,
		paymentsService: paymentsService,
		tokenManager:    tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initStudentsRoutes(v1)
		h.initCallbackRoutes(v1)
	}
}
