package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/service"
	"github.com/zhashkevych/courses-backend/pkg/auth"
)

type Handler struct {
	schoolsService    service.Schools
	studentsService   service.Students
	coursesService    service.Courses
	promoCodesService service.PromoCodes
	offersService     service.Offers
	modulesService    service.Modules
	ordersService     service.Orders
	paymentsService   service.Payments
	adminsService     service.Admins
	packagesService   service.Packages
	lessonsService    service.Lessons
	tokenManager      auth.TokenManager
}

func NewHandler(schoolsService service.Schools, studentsService service.Students, coursesService service.Courses, promoCodesService service.PromoCodes,
	offersService service.Offers, modulesService service.Modules, ordersService service.Orders,
	paymentsService service.Payments, adminsService service.Admins, packagesService service.Packages, lessonsService service.Lessons, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		schoolsService:    schoolsService,
		studentsService:   studentsService,
		coursesService:    coursesService,
		offersService:     offersService,
		promoCodesService: promoCodesService,
		modulesService:    modulesService,
		ordersService:     ordersService,
		paymentsService:   paymentsService,
		adminsService:     adminsService,
		packagesService:   packagesService,
		lessonsService:    lessonsService,
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
	}
}
