package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: review response error messages

func (h *Handler) initAdminRoutes(api *gin.RouterGroup) { //nolint:funlen
	admins := api.Group("/admins", h.setSchoolFromRequest)
	{
		admins.POST("/sign-in", h.adminSignIn)
		admins.POST("/auth/refresh", h.adminRefresh)

		authenticated := admins.Group("/", h.adminIdentity)
		{
			courses := authenticated.Group("/courses")
			{
				courses.POST("", h.adminCreateCourse)
				courses.GET("", h.adminGetAllCourses)
				courses.GET("/:id", h.adminGetCourseById)
				courses.PUT("/:id", h.adminUpdateCourse)
				courses.DELETE("/:id", h.adminDeleteCourse)
				courses.POST("/:id/modules", h.adminCreateModule)
				courses.POST("/:id/packages", h.adminCreatePackage)
				courses.GET("/:id/packages", h.adminGetAllPackages)
			}

			modules := authenticated.Group("/modules")
			{
				modules.PUT("/:id", h.adminUpdateModule)
				modules.DELETE("/:id", h.adminDeleteModule)
				modules.GET("/:id/lessons", h.adminGetLessons)
				modules.POST("/:id/lessons", h.adminCreateLesson)
			}

			lessons := authenticated.Group("/lessons")
			{
				lessons.GET("/:id", h.adminGetLessonById)
				lessons.PUT("/:id", h.adminUpdateLesson)
				lessons.DELETE("/:id", h.adminDeleteLesson)
			}

			packages := authenticated.Group("/packages")
			{
				packages.GET("/:id", h.adminGetPackageById)
				packages.PUT("/:id", h.adminUpdatePackage)
				packages.DELETE("/:id", h.adminDeletePackage)
			}

			offers := authenticated.Group("/offers")
			{
				offers.POST("", h.adminCreateOffer)
				offers.GET("", h.adminGetAllOffers)
				offers.GET("/:id", h.adminGetOfferById)
				offers.PUT("/:id", h.adminUpdateOffer)
				offers.DELETE("/:id", h.adminDeleteOffer)
			}

			school := authenticated.Group("/school")
			{
				school.PUT("/settings", h.adminUpdateSchoolSettings)
			}

			promocodes := authenticated.Group("/promocodes")
			{
				promocodes.POST("", h.adminCreatePromocode)
				promocodes.GET("", h.adminGetPromocodes)
				promocodes.GET("/:id", h.adminGetPromocodeById)
				promocodes.PUT("/:id", h.adminUpdatePromocode)
				promocodes.DELETE("/:id", h.adminDeletePromocode)
			}

			authenticated.GET("/orders", h.adminGetOrders)
			authenticated.GET("/students", h.adminGetStudents)

			upload := authenticated.Group("/upload")
			{
				upload.POST("/image", h.adminUploadImage)
				upload.POST("/video", h.adminUploadVideo)
			}
		}
	}
}

// @Summary Admin SignIn
// @Tags admins-auth
// @Description admin sign in
// @ModuleID adminSignIn
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

	res, err := h.services.Admins.SignIn(c.Request.Context(), service.SchoolSignInInput{
		Email:    inp.Email,
		Password: inp.Password,
		SchoolID: school.ID,
	})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			newResponse(c, http.StatusUnauthorized, err.Error())
		} else {
			newResponse(c, http.StatusInternalServerError, err.Error())
		}

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

	res, err := h.services.Admins.RefreshTokens(c.Request.Context(), school.ID, inp.Token)
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
	Name string `json:"name" binding:"required"`
}

// @Summary Admin Create New Courses
// @Security AdminAuth
// @Tags admins-courses
// @Description admin create new course
// @ModuleID adminCreateCourse
// @Accept  json
// @Produce  json
// @Param input body createCourseInput true "course info"
// @Success 200 {object} idResponse
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

	id, err := h.services.Courses.Create(c.Request.Context(), school.ID, inp.Name)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, idResponse{id})
}

// @Summary Admin Get All Courses
// @Security AdminAuth
// @Tags admins-courses
// @Description admin get all courses
// @ModuleID adminGetAllCourses
// @Accept  json
// @Produce  json
// @Success 200 {object} dataResponse
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

	courses, err := h.services.Admins.GetCourses(c.Request.Context(), school.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{courses})
}

type adminGetCourseByIdResponse struct {
	Course  domain.Course   `json:"course"`
	Modules []domain.Module `json:"modules"`
}

// @Summary Admin Get Course By ID
// @Security AdminAuth
// @Tags admins-courses
// @Description admin get course by id
// @ModuleID adminGetCourseById
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} domain.Course
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id} [get]
func (h *Handler) adminGetCourseById(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	course, err := h.services.Admins.GetCourseById(c.Request.Context(), school.ID, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	modules, err := h.services.Modules.GetByCourseId(c.Request.Context(), course.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, adminGetCourseByIdResponse{
		Course:  course,
		Modules: modules,
	})
}

type updateCourseInput struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Published   *bool  `json:"published"`
}

// @Summary Admin Update Course
// @Security AdminAuth
// @Tags admins-courses
// @Description admin update course
// @ModuleID adminUpdateCourse
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Param input body updateCourseInput true "course update info"
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

	var inp updateCourseInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Courses.Update(c.Request.Context(), service.UpdateCourseInput{
		CourseID:    idParam,
		SchoolID:    school.ID.Hex(),
		Name:        inp.Name,
		Description: inp.Description,
		Code:        inp.Code,
		Color:       inp.Color,
		Published:   inp.Published,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// TODO cover with tests
// @Summary Admin Delete Course
// @Security AdminAuth
// @Tags admins-courses
// @Description admin delete course
// @ModuleID adminDeleteCourse
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id} [delete]
func (h *Handler) adminDeleteCourse(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Courses.Delete(c.Request.Context(), school.ID, id); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

type createModuleInput struct {
	Name     string `json:"name" binding:"required,min=5"`
	Position uint   `json:"position"`
}

// @Summary Admin Create Module
// @Security AdminAuth
// @Tags admins-modules
// @Description admin update course
// @ModuleID adminCreateModule
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Param input body createModuleInput true "module info"
// @Success 201 {object} idResponse
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

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	moduleId, err := h.services.Modules.Create(c.Request.Context(), service.CreateModuleInput{
		SchoolID: school.ID.Hex(),
		CourseID: id,
		Name:     inp.Name,
		Position: inp.Position,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "invalid id param")

		return
	}

	c.JSON(http.StatusCreated, idResponse{moduleId})
}

type updateModuleInput struct {
	Name      string `json:"name"`
	Position  *uint  `json:"position"`
	Published *bool  `json:"published"`
}

// @Summary Admin Update Module
// @Security AdminAuth
// @Tags admins-modules
// @Description admin update course
// @ModuleID adminUpdateModule
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

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Modules.Update(c.Request.Context(), service.UpdateModuleInput{
		ID:        id,
		SchoolID:  school.ID.Hex(),
		Name:      inp.Name,
		Position:  inp.Position,
		Published: inp.Published,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// @Summary Admin Delete Module
// @Security AdminAuth
// @Tags admins-modules
// @Description admin update course
// @ModuleID adminDeleteModule
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

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.Modules.Delete(c.Request.Context(), school.ID, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// @Summary Admin Get Module Lessons
// @Security AdminAuth
// @Tags admins-lessons
// @Description admin get module lessons with content
// @ModuleID adminGetLessons
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/modules/{id}/lessons [get]
func (h *Handler) adminGetLessons(c *gin.Context) {
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

	module, err := h.services.Modules.GetWithContent(c.Request.Context(), moduleId)
	if err != nil {
		if errors.Is(err, domain.ErrModuleIsNotAvailable) {
			newResponse(c, http.StatusForbidden, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{module.Lessons})
}

type createLessonInput struct {
	Name     string `json:"name" binding:"required,min=3"`
	Position uint   `json:"position"`
}

// @Summary Admin Create Lesson
// @Security AdminAuth
// @Tags admins-lessons
// @Description admin create lesson
// @ModuleID adminCreateLesson
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Param input body createLessonInput true "lesson info"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/modules/{id}/lessons [post]
func (h *Handler) adminCreateLesson(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	var inp createLessonInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	lessonId, err := h.services.Lessons.Create(c.Request.Context(), service.AddLessonInput{
		ModuleID: id,
		Name:     inp.Name,
		Position: inp.Position,
		SchoolID: school.ID.Hex(),
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, idResponse{lessonId})
}

// @Summary Admin Get Lesson By Id
// @Security AdminAuth
// @Tags admins-lessons
// @Description admin get lesson by Id
// @ModuleID adminGetLessonById
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/lessons/{id} [get]
func (h *Handler) adminGetLessonById(c *gin.Context) {
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

	lesson, err := h.services.Lessons.GetById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, lesson)
}

type updateLessonInput struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	Position  *uint  `json:"position"`
	Published *bool  `json:"published"`
}

// @Summary Admin Update Lesson
// @Security AdminAuth
// @Tags admins-lessons
// @Description admin update lesson
// @ModuleID adminUpdateLesson
// @Accept  json
// @Produce  json
// @Param id path string true "lesson id"
// @Param input body updateLessonInput true "update info"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/lessons/{id} [put]
func (h *Handler) adminUpdateLesson(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	var inp updateLessonInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Lessons.Update(c.Request.Context(), service.UpdateLessonInput{
		LessonID:  id,
		Name:      inp.Name,
		Content:   inp.Content,
		Position:  inp.Position,
		Published: inp.Published,
		SchoolID:  school.ID.Hex(),
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// @Summary Admin Delete Lesson
// @Security AdminAuth
// @Tags admins-lessons
// @Description admin delete lesson
// @ModuleID adminDeleteLesson
// @Accept  json
// @Produce  json
// @Param id path string true "lesson id"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/lessons/{id} [delete]
func (h *Handler) adminDeleteLesson(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.Lessons.Delete(c.Request.Context(), school.ID, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

type createPackageInput struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description"`
}

// @Summary Admin Create Package
// @Security AdminAuth
// @Tags admins-packages
// @Description admin create package
// @ModuleID adminCreatePackage
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Param input body createPackageInput true "package info"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id}/packages [post]
func (h *Handler) adminCreatePackage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	var inp createPackageInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	moduleId, err := h.services.Packages.Create(c.Request.Context(), service.CreatePackageInput{
		SchoolID:    school.ID.Hex(),
		CourseID:    id,
		Name:        inp.Name,
		Description: inp.Description,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, idResponse{moduleId})
}

// @Summary Admin Get All Course Packages
// @Security AdminAuth
// @Tags admins-packages
// @Description admin get all course packages
// @ModuleID adminGetAllPackages
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id}/packages [get]
func (h *Handler) adminGetAllPackages(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	packages, err := h.services.Packages.GetByCourse(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "invalid id param")

		return
	}

	c.JSON(http.StatusOK, dataResponse{packages})
}

// @Summary Admin Get Package By ID
// @Security AdminAuth
// @Tags admins-packages
// @Description admin get package by id
// @ModuleID adminGetPackageById
// @Accept  json
// @Produce  json
// @Param id path string true "package id"
// @Success 200 {array} domain.Package
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/packages/{id} [get]
func (h *Handler) adminGetPackageById(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	pkg, err := h.services.Packages.GetById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "invalid id param")

		return
	}

	c.JSON(http.StatusOK, pkg)
}

type updatePackageInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Modules     []string `json:"modules"`
}

// @Summary Admin Update Package
// @Security AdminAuth
// @Tags admins-packages
// @Description admin update package
// @ModuleID adminUpdatePackage
// @Accept  json
// @Produce  json
// @Param id path string true "package id"
// @Param input body updatePackageInput true "update input"
// @Success 200 {array} domain.Package
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/packages/{id} [put]
func (h *Handler) adminUpdatePackage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	var inp updatePackageInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Packages.Update(c.Request.Context(), service.UpdatePackageInput{
		ID:          id,
		SchoolID:    school.ID.Hex(),
		Name:        inp.Name,
		Description: inp.Description,
		Modules:     inp.Modules,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// @Summary Admin Delete Package
// @Security AdminAuth
// @Tags admins-packages
// @Description admin delete package
// @ModuleID adminDeletePackage
// @Accept  json
// @Produce  json
// @Param id path string true "package id"
// @Success 200 {array} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/packages/{id} [delete]
func (h *Handler) adminDeletePackage(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.Packages.Delete(c.Request.Context(), school.ID, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "invalid id param")

		return
	}

	c.Status(http.StatusOK)
}

type createOfferInput struct {
	Name        string   `json:"name" binding:"required,min=3"`
	Description string   `json:"description"`
	Benefits    []string `json:"benefits" binding:"required"`
	Price       price    `json:"price" binding:"required"`
}

// @Summary Admin Create Offer
// @Security AdminAuth
// @Tags admins-offers
// @Description admin create offer
// @ModuleID adminCreateOffer
// @Accept  json
// @Produce  json
// @Param input body createOfferInput true "package info"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/offers [post]
func (h *Handler) adminCreateOffer(c *gin.Context) {
	var inp createOfferInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	id, err := h.services.Offers.Create(c.Request.Context(), service.CreateOfferInput{
		SchoolID:    school.ID,
		Name:        inp.Name,
		Description: inp.Description,
		Benefits:    inp.Benefits,
		Price: domain.Price{
			Value:    inp.Price.Value,
			Currency: inp.Price.Currency,
		},
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, idResponse{id})
}

// @Summary Admin Get All Offers
// @Security AdminAuth
// @Tags admins-offers
// @Description admin get all offers
// @ModuleID adminGetAllOffers
// @Accept  json
// @Produce  json
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/offers [get]
func (h *Handler) adminGetAllOffers(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	offers, err := h.services.Offers.GetAll(c.Request.Context(), school.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{offers})
}

// @Summary Admin Get Offer By Id
// @Security AdminAuth
// @Tags admins-offers
// @Description admin get offer by id
// @ModuleID adminGetOfferById
// @Accept  json
// @Produce  json
// @Param id path string true "offer id"
// @Success 200 {object} domain.Offer
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/offers/{id} [get]
func (h *Handler) adminGetOfferById(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	offer, err := h.services.Offers.GetById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, offer)
}

type updateOfferInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Benefits    []string `json:"benefits"`
	Price       *price   `json:"price"`
	Packages    []string `json:"packages"`
}

// @Summary Admin Update Offer
// @Security AdminAuth
// @Tags admins-offers
// @Description admin updateOffer
// @ModuleID adminUpdateOffer
// @Accept  json
// @Produce  json
// @Param id path string true "offer id"
// @Param input body updateOfferInput true "update info"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/offers/{id} [put]
func (h *Handler) adminUpdateOffer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	var inp updateOfferInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	updateInput := service.UpdateOfferInput{
		ID:          id,
		SchoolID:    school.ID.Hex(),
		Name:        inp.Name,
		Description: inp.Description,
		Packages:    inp.Packages,
		Benefits:    inp.Benefits,
	}

	if inp.Price != nil {
		updateInput.Price = &domain.Price{
			Value:    inp.Price.Value,
			Currency: inp.Price.Currency,
		}
	}

	if err := h.services.Offers.Update(c.Request.Context(), updateInput); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// @Summary Admin Delete Offer
// @Security AdminAuth
// @Tags admins-offers
// @Description admin delete offer
// @ModuleID adminDeleteOffer
// @Accept  json
// @Produce  json
// @Param id path string true "offer id"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/offers/{id} [delete]
func (h *Handler) adminDeleteOffer(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.Offers.Delete(c.Request.Context(), school.ID, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

type createPromocodeInput struct {
	Code               string               `json:"code" binding:"required"`
	DiscountPercentage int                  `json:"discountPercentage" binding:"required"`
	ExpiresAt          time.Time            `json:"expiresAt" binding:"required"`
	OfferIDs           []primitive.ObjectID `json:"offerIds" binding:"required"`
}

// @Summary Admin Create Promocode
// @Security AdminAuth
// @Tags admins-promocodes
// @Description admin create promocode
// @ModuleID adminCreatePromocode
// @Accept  json
// @Produce  json
// @Param input body createPromocodeInput true "package info"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/promocodes [post]
func (h *Handler) adminCreatePromocode(c *gin.Context) {
	var inp createPromocodeInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	id, err := h.services.PromoCodes.Create(c.Request.Context(), service.CreatePromoCodeInput{
		SchoolID:           school.ID,
		Code:               inp.Code,
		DiscountPercentage: inp.DiscountPercentage,
		ExpiresAt:          inp.ExpiresAt,
		OfferIDs:           inp.OfferIDs,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, idResponse{id})
}

// @Summary Admin Get All Promocodes
// @Security AdminAuth
// @Tags admins-promocodes
// @Description admin get all promocodes
// @ModuleID adminGetPromocodes
// @Accept  json
// @Produce  json
// @Success 200 {object} dataResponse{data=[]domain.PromoCode}
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/promocodes [get]
func (h *Handler) adminGetPromocodes(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	promocodes, err := h.services.PromoCodes.GetBySchool(c.Request.Context(), school.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{promocodes})
}

// @Summary Admin Get Promocode By Id
// @Security AdminAuth
// @Tags admins-promocodes
// @Description admin get promocode by id
// @ModuleID adminGetPromocodeById
// @Accept  json
// @Produce  json
// @Param id path string true "promocode id"
// @Success 200 {object} domain.PromoCode
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/promocodes/{id} [get]
func (h *Handler) adminGetPromocodeById(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	promoCode, err := h.services.PromoCodes.GetById(c.Request.Context(), school.ID, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, promoCode)
}

type updatePromocodeInput struct {
	Code               string    `json:"code"`
	DiscountPercentage int       `json:"discountPercentage"`
	ExpiresAt          time.Time `json:"expiresAt"`
	OfferIDs           []string  `json:"offerIds"`
}

// @Summary Admin Update Promocode
// @Security AdminAuth
// @Tags admins-promocodes
// @Description admin update promocode
// @ModuleID adminUpdatePromocode
// @Accept  json
// @Produce  json
// @Param id path string true "promocode id"
// @Param input body updatePromocodeInput true "update info"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/promocodes/{id} [put]
func (h *Handler) adminUpdatePromocode(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	var inp updatePromocodeInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.PromoCodes.Update(c.Request.Context(), service.UpdatePromoCodeInput{
		ID:                 id,
		SchoolID:           school.ID,
		Code:               inp.Code,
		DiscountPercentage: inp.DiscountPercentage,
		ExpiresAt:          inp.ExpiresAt,
		OfferIDs:           inp.OfferIDs,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// @Summary Admin Delete Promocode
// @Security AdminAuth
// @Tags admins-promocodes
// @Description admin delete promocode
// @ModuleID adminDeletePromocode
// @Accept  json
// @Produce  json
// @Param id path string true "promocode id"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/promocodes/{id} [delete]
func (h *Handler) adminDeletePromocode(c *gin.Context) {
	id, err := parseIdFromPath(c)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	err = h.services.PromoCodes.Delete(c.Request.Context(), school.ID, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

type pages struct {
	Confidential     string `json:"confidential"`
	ServiceAgreement string `json:"serviceAgreement"`
	RefundPolicy     string `json:"refundPolicy"`
}

type contactInfo struct {
	BusinessName       string `json:"businessName"`
	RegistrationNumber string `json:"registrationNumber"`
	Address            string `json:"address"`
	Email              string `json:"email"`
}

type updateSchoolSettingsInput struct {
	Color       string       `json:"color"`
	Domains     []string     `json:"domains"`
	Email       string       `json:"email"`
	ContactInfo *contactInfo `json:"contactInfo"`
	Pages       *pages       `json:"pages"`
}

// @Summary Admin Update School settings
// @Security AdminAuth
// @Tags admins-school
// @Description admin update school settings
// @ModuleID adminUpdateSchoolSettings
// @Accept  json
// @Produce  json
// @Param input body updateSchoolSettingsInput true "update school settings"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/school/settings [put]
func (h *Handler) adminUpdateSchoolSettings(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	var inp updateSchoolSettingsInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	updateInput := service.UpdateSchoolSettingsInput{
		Color:   inp.Color,
		Domains: inp.Domains,
		Email:   inp.Email,
	}

	if inp.Pages != nil {
		updateInput.Pages = &domain.Pages{
			Confidential:     inp.Pages.Confidential,
			ServiceAgreement: inp.Pages.ServiceAgreement,
			RefundPolicy:     inp.Pages.RefundPolicy,
		}
	}

	if inp.ContactInfo != nil {
		updateInput.ContactInfo = &domain.ContactInfo{
			Email:              inp.ContactInfo.Email,
			RegistrationNumber: inp.ContactInfo.RegistrationNumber,
			Address:            inp.ContactInfo.Address,
			BusinessName:       inp.ContactInfo.BusinessName,
		}
	}

	if err := h.services.Schools.UpdateSettings(c.Request.Context(), school.ID, updateInput); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// @Summary Admin Get Orders
// @Security AdminAuth
// @Tags admins-orders
// @Description admin get all orders
// @ModuleID adminGetOrders
// @Accept  json
// @Produce  json
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/orders [get]
func (h *Handler) adminGetOrders(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	orders, err := h.services.Orders.GetBySchool(c.Request.Context(), school.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{orders})
}

type studentResponse struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	RegisteredAt time.Time          `json:"registeredAt"`
	LastVisitAt  time.Time          `json:"lastVisitAt"`
}

func toStudentsResponse(students []domain.Student) []studentResponse {
	out := make([]studentResponse, len(students))
	for i, student := range students {
		out[i].ID = student.ID
		out[i].Name = student.Name
		out[i].Email = student.Email
		out[i].RegisteredAt = student.RegisteredAt
		out[i].LastVisitAt = student.LastVisitAt
	}

	return out
}

// @Summary Admin Get Students
// @Security AdminAuth
// @Tags admins-orders
// @Description admin get all students
// @ModuleID adminGetStudents
// @Accept  json
// @Produce  json
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/students [get]
func (h *Handler) adminGetStudents(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	students, err := h.services.Students.GetBySchool(c.Request.Context(), school.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{toStudentsResponse(students)})
}
