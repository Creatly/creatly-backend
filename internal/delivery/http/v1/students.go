package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/service"
	"github.com/zhashkevych/courses-backend/pkg/logger"
	"net/http"
)

func (h *Handler) initStudentsRoutes(api *gin.RouterGroup) {
	students := api.Group("/students", h.setSchoolFromRequest())
	{
		students.POST("/sign-up", h.studentSignUp)
		students.POST("/sign-in")
		students.POST("/verify/:hash")
		students.POST("/courses")
		students.POST("/courses/:id")
	}
}

type studentSignUpInput struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	SourceCourseId string `json:"sourceCourseId"`
}

func (h *Handler) studentSignUp(c *gin.Context) {
	var inp studentSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.studentsService.SignUp(c.Request.Context(), service.StudentSignUpInput{
		Name:           inp.Name,
		Email:          inp.Email,
		Password:       inp.Password,
		SourceCourseID: inp.SourceCourseId,
		SchoolID:       school.ID,
	}); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
