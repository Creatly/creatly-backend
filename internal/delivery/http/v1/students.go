package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/service"
	"net/http"
)

func (h *Handler) initStudentsRoutes(api *gin.RouterGroup) {
	students := api.Group("/students", h.setSchoolFromRequest)
	{
		students.POST("/sign-up", h.studentSignUp)
		students.POST("/sign-in", h.studentSignIn)
		students.POST("/auth/refresh", h.studentRefresh)
		students.POST("/verify/:hash", h.studentVerify)

		authenticated := students.Group("/", h.userIdentity)
		{
			authenticated.GET("/courses", h.studentGetAllCourses)
			authenticated.GET("/courses/:id", h.studentGetCourseById)
		}
	}
}

type studentSignUpInput struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RegisterSource string `json:"registerSource"`
}

func (h *Handler) studentSignUp(c *gin.Context) {
	var inp studentSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.studentsService.SignUp(c.Request.Context(), service.StudentSignUpInput{
		Name:           inp.Name,
		Email:          inp.Email,
		Password:       inp.Password,
		RegisterSource: inp.RegisterSource,
		SchoolID:       school.ID,
	}); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

type studentSignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) studentSignIn(c *gin.Context) {
	var inp studentSignInInput
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := h.studentsService.SignIn(c.Request.Context(), service.StudentSignInInput{
		SchoolID: school.ID,
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

func (h *Handler) studentRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := h.studentsService.RefreshTokens(c.Request.Context(), school.ID, inp.Token)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) studentVerify(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newErrorResponse(c, http.StatusBadRequest, "code is empty")
		return
	}

	if err := h.studentsService.Verify(c.Request.Context(), code); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) studentGetAllCourses(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Return only published courses
	courses := make([]domain.Course, 0)
	for _, course := range school.Courses {
		if course.Published {
			courses = append(courses, course)
		}
	}

	c.JSON(http.StatusOK, courses)
}

func (h *Handler) studentGetCourseById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for _, course := range school.Courses {
		if course.Published && course.ID.Hex() == id {
			c.JSON(http.StatusOK, course)
			return
		}
	}

	newErrorResponse(c, http.StatusBadRequest, "not found")
}
