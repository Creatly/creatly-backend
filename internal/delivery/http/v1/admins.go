package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *Handler) initAdminRoutes(api *gin.RouterGroup) {
	students := api.Group("/admins", h.setSchoolFromRequest)
	{
		students.POST("/sign-in", h.adminSignIn)
		students.POST("/auth/refresh", h.adminRefresh)

		authenticated := students.Group("/", h.adminIdentity)
		{
			courses := authenticated.Group("/courses")
			{
				courses.POST("/", h.adminCreateCourse)
				courses.GET("/", h.adminGetAllCourses)
				courses.GET("/:id", h.adminGetCourseById)
				courses.PUT("/:id", h.adminUpdateCourse)
				courses.POST("/:id/modules", h.adminCreateModule)
			}

			modules := authenticated.Group("/modules")
			{
				modules.PUT("/:id", h.adminUpdateModule)
				modules.DELETE("/:id", h.adminDeleteModule)
				modules.GET(":id/lessons", h.adminGetModuleLessons)
				modules.POST("id/lessons", h.adminCreateLesson)
			}

		}
	}
}

// @Summary Admin SignIn
// @Tags admins-auth
// @Description admin sign in
// @ID adminSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/sign-in [post]
func (h *Handler) adminSignIn(c *gin.Context) {
	var inp signInInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := h.adminsService.SignIn(c.Request.Context(), service.SignInInput{
		Email:    inp.Email,
		Password: inp.Password,
		SchoolID: school.ID,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// @Summary Admin Refresh Tokens
// @Tags admins-auth
// @Description admin refresh tokens
// @Accept  json
// @Produce  json
// @Param input body refreshInput true "refresh info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/auth/refresh [post]
func (h *Handler) adminRefresh(c *gin.Context) {
	var inp refreshInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := h.adminsService.RefreshTokens(c.Request.Context(), school.ID, inp.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

type createCourseInput struct {
	Name string `json:"name,required"`
}

// @Summary Admin Create New Courses
// @Security AdminAuth
// @Tags admins-courses
// @Description admin create new course
// @ID adminCreateCourse
// @Accept  json
// @Produce  json
// @Param input body createCourseInput true "course info"
// @Success 200 {array} domain.Course
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses [post]
func (h *Handler) adminCreateCourse(c *gin.Context) {
	var inp createCourseInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.coursesService.Create(c.Request.Context(), school.ID, inp.Name)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// @Summary Admin Get All Courses
// @Security AdminAuth
// @Tags admins-courses
// @Description admin get all courses
// @ID adminGetAllCourses
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Course
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses [get]
func (h *Handler) adminGetAllCourses(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	courses, err := h.adminsService.GetCourses(c.Request.Context(), school.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, courses)
}

type adminGetCourseByIdResponse struct {
	Course  domain.Course   `json:"course"`
	Modules []domain.Module `json:"modules"`
}

// @Summary Admin Get Course By ID
// @Security AdminAuth
// @Tags admins-courses
// @Description admin get course by id
// @ID adminGetCourseById
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} domain.Course
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id} [get]
func (h *Handler) adminGetCourseById(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	course, err := h.adminsService.GetCourseById(c.Request.Context(), school.ID, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	modules, err := h.modulesService.GetByCourse(c.Request.Context(), course.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, adminGetCourseByIdResponse{
		Course:  course,
		Modules: modules,
	})
}

type adminUpdateCourseInput struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Published   *bool  `json:"published"`
}

// @Summary Admin Update Course
// @Security AdminAuth
// @Tags admins-courses
// @Description admin update course
// @ID adminUpdateCourse
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Param input body adminUpdateCourseInput true "course update info"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id} [put]
func (h *Handler) adminUpdateCourse(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	var inp adminUpdateCourseInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.coursesService.Update(c.Request.Context(), school.ID, service.UpdateCourseInput{
		CourseID:    idParam,
		Name:        inp.Name,
		Description: inp.Description,
		Code:        inp.Code,
		Published:   inp.Published,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Admin Get Module Lessons
// @Security AdminAuth
// @Tags admins-lessons
// @Description admin get module content
// @ID adminGetModuleLessons
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/modules/{id}/lessons [get]
func (h *Handler) adminGetModuleLessons(c *gin.Context) {
	moduleIdParam := c.Param("id")
	if moduleIdParam == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	moduleId, err := primitive.ObjectIDFromHex(moduleIdParam)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	module, err := h.modulesService.GetWithContent(c.Request.Context(), moduleId)
	if err != nil {
		if err == service.ErrModuleIsNotAvailable {
			newResponse(c, http.StatusForbidden, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getModuleLessonsResponse{
		Lessons: module.Lessons,
	})
}

type createModuleInput struct {
	Name     string `json:"name" binding:"required,min=5"`
	Position int    `json:"position" binding:"required,min=0"`
}

// @Summary Admin Create Module
// @Security AdminAuth
// @Tags admins-modules
// @Description admin update course
// @ID adminCreateModule
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Param input body createModuleInput true "module info"
// @Success 201 {string} string "id"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id}/modules [post]
func (h *Handler) adminCreateModule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	var inp createModuleInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	moduleId, err := h.modulesService.Create(c.Request.Context(), service.CreateModuleInput{
		CourseID: id,
		Name:     inp.Name,
		Position: inp.Position,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": moduleId,
	})
}

type updateModuleInput struct {
	Name      string `json:"name"`
	Position  *int   `json:"position"`
	Published *bool  `json:"published"`
}

// @Summary Admin Update Module
// @Security AdminAuth
// @Tags admins-modules
// @Description admin update course
// @ID adminUpdateModule
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Param input body updateModuleInput true "update info"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/modules/{id} [put]
func (h *Handler) adminUpdateModule(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	var inp updateModuleInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.modulesService.Update(c.Request.Context(), service.UpdateModuleInput{
		ID:        id,
		Name:      inp.Name,
		Position:  inp.Position,
		Published: inp.Published,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Admin Delete Module
// @Security AdminAuth
// @Tags admins-modules
// @Description admin update course
// @ID adminDeleteModule
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/modules/{id} [delete]
func (h *Handler) adminDeleteModule(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")
		return
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.modulesService.Delete(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) adminCreateLesson(c *gin.Context) {

}

func (h *Handler) adminUpdateLesson(c *gin.Context) {

}

func (h *Handler) adminDeleteLesson(c *gin.Context) {

}
