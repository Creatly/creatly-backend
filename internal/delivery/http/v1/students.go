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
		students.POST("/sign-in", h.studentSignIn)
		students.POST("/auth/refresh", h.studentRefresh)
		students.POST("/verify/:hash", h.studentVerify)
		students.POST("/courses")
		students.POST("/courses/:id")
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
		RegisterSource: inp.RegisterSource,
		SchoolID:       school.ID,
	}); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res, err := h.studentsService.SignIn(c.Request.Context(), service.StudentSignInInput{
		SchoolID: school.ID,
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	res, err := h.studentsService.RefreshTokens(c.Request.Context(), school.ID, inp.Token)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *Handler) studentVerify(c *gin.Context) {
	hash := c.Param("hash")
	if hash == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.studentsService.Verify(c.Request.Context(), hash); err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}
