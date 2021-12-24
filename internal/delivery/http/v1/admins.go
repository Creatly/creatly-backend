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

				modules.GET("/:id/survey", h.adminGetSurvey)
				modules.POST("/:id/survey", h.adminCreateOrUpdateSurvey)
				modules.DELETE("/:id/survey", h.adminDeleteSurvey)
				modules.GET("/:id/survey/results", h.adminGetSurveyResults)
				modules.GET("/:id/survey/results/:studentId", h.adminGetSurveyStudentResults)
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
				school.PUT("/settings/fondy", h.adminConnectFondy)
				school.PUT("/settings/sendpulse", h.adminConnectSendPulse)
			}

			promocodes := authenticated.Group("/promocodes")
			{
				promocodes.POST("", h.adminCreatePromocode)
				promocodes.GET("", h.adminGetPromocodes)
				promocodes.GET("/:id", h.adminGetPromocodeById)
				promocodes.PUT("/:id", h.adminUpdatePromocode)
				promocodes.DELETE("/:id", h.adminDeletePromocode)
			}

			orders := authenticated.Group("/orders")
			{
				orders.GET("", h.adminGetOrders)
				orders.PUT("/:id", h.adminUpdateOrderStatus)
			}

			students := authenticated.Group("/students")
			{
				students.GET("", h.adminGetStudents)
				students.POST("", h.adminCreateStudent)
				students.GET("/:id", h.adminGetStudentById)
				students.PUT("/:id", h.adminUpdateStudent)
				students.DELETE("/:id", h.adminDeleteStudent)
				students.PATCH("/:id/offers/:offerId", h.adminManageOfferPermission)
			}

			upload := authenticated.Group("/upload")
			{
				upload.POST("/image", h.adminUploadImage)
				upload.POST("/video", h.adminUploadVideo)
				upload.POST("/file", h.adminUploadFile)
			}

			media := authenticated.Group("/media")
			{
				media.GET("/videos/:id", h.adminGetVideo)
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

	response := make([]domain.Course, len(courses))
	if courses != nil {
		response = courses
	}

	c.JSON(http.StatusOK, dataResponse{Data: response})
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
	id, err := parseIdFromPath(c, "id")
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
	Name        *string `json:"name"`
	ImageURL    *string `json:"imageUrl"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
	Published   *bool   `json:"published"`
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
		ImageURL:    inp.ImageURL,
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
	id, err := parseIdFromPath(c, "id")
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

	c.JSON(http.StatusOK, dataResponse{Data: module.Lessons})
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
	id, err := parseIdFromPath(c, "id")
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
	Name    string   `json:"name" binding:"required,min=3"`
	Modules []string `json:"modules"`
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
		SchoolID: school.ID.Hex(),
		CourseID: id,
		Name:     inp.Name,
		Modules:  inp.Modules,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, idResponse{moduleId})
}

type packageResponse struct {
	ID      primitive.ObjectID `json:"id"`
	Name    string             `json:"name"`
	Modules []packageModule    `json:"modules,omitempty"`
}

type packageModule struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
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
	id, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	pkg, err := h.services.Packages.GetByCourse(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "invalid id param")

		return
	}

	c.JSON(http.StatusOK, dataResponse{Data: toPackagesResponse(pkg)})
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
	id, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	pkg, err := h.services.Packages.GetById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "invalid id param")

		return
	}

	c.JSON(http.StatusOK, packageResponse{
		ID:      pkg.ID,
		Name:    pkg.Name,
		Modules: toPackageModules(pkg.Modules),
	})
}

type updatePackageInput struct {
	Name    string   `json:"name"`
	Modules []string `json:"modules"`
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
		ID:       id,
		SchoolID: school.ID.Hex(),
		Name:     inp.Name,
		Modules:  inp.Modules,
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
	id, err := parseIdFromPath(c, "id")
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
	Name          string        `json:"name" binding:"required,min=3"`
	Description   string        `json:"description"`
	Benefits      []string      `json:"benefits" binding:"required"`
	Packages      []string      `json:"packages"`
	Price         price         `json:"price" binding:"required"`
	PaymentMethod paymentMethod `json:"paymentMethod" binding:"required"`
}

type paymentMethod struct {
	UsesProvider bool   `json:"usesProvider"`
	Provider     string `json:"provider"`
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
		PaymentMethod: domain.PaymentMethod{
			UsesProvider: inp.PaymentMethod.UsesProvider,
			Provider:     inp.PaymentMethod.Provider,
		},
		Packages: inp.Packages,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusCreated, idResponse{id})
}

type offerResponse struct {
	ID            primitive.ObjectID   `json:"id"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	Benefits      []string             `json:"benefits"`
	Packages      []packageResponse    `json:"packages"`
	Price         domain.Price         `json:"price"`
	PaymentMethod domain.PaymentMethod `json:"paymentMethod"`
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

	response := make([]offerResponse, len(offers))

	for i, offer := range offers {
		pkgs, err := h.services.Packages.GetByIds(c.Request.Context(), offer.PackageIDs)
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())

			return
		}

		response[i] = offerResponse{
			ID:            offer.ID,
			Name:          offer.Name,
			Description:   offer.Description,
			Benefits:      offer.Benefits,
			Price:         offer.Price,
			PaymentMethod: offer.PaymentMethod,
			Packages:      toPackagesResponse(pkgs),
		}
	}

	c.JSON(http.StatusOK, dataResponse{Data: response})
}

// @Summary Admin Get Offer By Id
// @Security AdminAuth
// @Tags admins-offers
// @Description admin get offer by id
// @ModuleID adminGetOfferById
// @Accept  json
// @Produce  json
// @Param id path string true "offer id"
// @Success 200 {object} offerResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/offers/{id} [get]
func (h *Handler) adminGetOfferById(c *gin.Context) {
	id, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	offer, err := h.services.Offers.GetById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	pkgs, err := h.services.Packages.GetByIds(c.Request.Context(), offer.PackageIDs)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	response := offerResponse{
		ID:            offer.ID,
		Name:          offer.Name,
		Description:   offer.Description,
		Benefits:      offer.Benefits,
		Price:         offer.Price,
		PaymentMethod: offer.PaymentMethod,
		Packages:      toPackagesResponse(pkgs),
	}

	c.JSON(http.StatusOK, response)
}

type updateOfferInput struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Benefits      []string       `json:"benefits"`
	Price         *price         `json:"price"`
	Packages      []string       `json:"packages"`
	PaymentMethod *paymentMethod `json:"paymentMethod"`
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

	if inp.PaymentMethod != nil {
		updateInput.PaymentMethod = &domain.PaymentMethod{
			UsesProvider: inp.PaymentMethod.UsesProvider,
			Provider:     inp.PaymentMethod.Provider,
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
	id, err := parseIdFromPath(c, "id")
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

type promocodeReponse struct {
	ID                 primitive.ObjectID  `json:"id"`
	Code               string              `json:"code"`
	DiscountPercentage int                 `json:"discountPercentage"`
	ExpiresAt          time.Time           `json:"expiresAt"`
	Offers             []offerShortReponse `json:"offers"`
}

type offerShortReponse struct {
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
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

	response := make([]promocodeReponse, len(promocodes))

	for i, promocode := range promocodes {
		offers, err := h.services.Offers.GetByIds(c.Request.Context(), promocode.OfferIDs)
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())

			return
		}

		response[i] = promocodeReponse{
			ID:                 promocode.ID,
			Code:               promocode.Code,
			DiscountPercentage: promocode.DiscountPercentage,
			ExpiresAt:          promocode.ExpiresAt,
			Offers:             toOffersReponse(offers),
		}
	}

	c.JSON(http.StatusOK, dataResponse{Data: response})
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
	id, err := parseIdFromPath(c, "id")
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

	offers, err := h.services.Offers.GetByIds(c.Request.Context(), promoCode.OfferIDs)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	response := promocodeReponse{
		ID:                 promoCode.ID,
		Code:               promoCode.Code,
		DiscountPercentage: promoCode.DiscountPercentage,
		ExpiresAt:          promoCode.ExpiresAt,
		Offers:             toOffersReponse(offers),
	}

	c.JSON(http.StatusOK, response)
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
	id, err := parseIdFromPath(c, "id")
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

	var offerIds []primitive.ObjectID
	if inp.OfferIDs != nil {
		offerIds, err = stringArrayToObjectId(inp.OfferIDs)
		if err != nil {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}
	}

	if err := h.services.PromoCodes.Update(c.Request.Context(), domain.UpdatePromoCodeInput{
		ID:                 id,
		SchoolID:           school.ID,
		Code:               inp.Code,
		DiscountPercentage: inp.DiscountPercentage,
		ExpiresAt:          inp.ExpiresAt,
		OfferIDs:           offerIds,
	}); err != nil {
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
	id, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.PromoCodes.Delete(c.Request.Context(), school.ID, id); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

type (
	pages struct {
		Confidential      *string `json:"confidential"`
		ServiceAgreement  *string `json:"serviceAgreement"`
		NewsletterConsent *string `json:"newsletterConsent"`
	}

	contactInfo struct {
		BusinessName       *string `json:"businessName"`
		RegistrationNumber *string `json:"registrationNumber"`
		Address            *string `json:"address"`
		Email              *string `json:"email"`
		Phone              *string `json:"phone"`
	}

	updateSchoolSettingsInput struct {
		Name                *string      `json:"name"`
		Color               *string      `json:"color"`
		Domains             []string     `json:"domains"`
		Email               *string      `json:"email"`
		ContactInfo         *contactInfo `json:"contactInfo"`
		Pages               *pages       `json:"pages"`
		ShowPaymentImages   *bool        `json:"showPaymentImages"`
		GoogleAnalyticsCode *string      `json:"googleAnalyticsCode"`
		LogoURL             *string      `json:"logo"`
		DisableRegistration *bool        `json:"disableRegistration"`
	}
)

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

	updateInput := domain.UpdateSchoolSettingsInput{
		Name:                inp.Name,
		Color:               inp.Color,
		Domains:             inp.Domains,
		Email:               inp.Email,
		ShowPaymentImages:   inp.ShowPaymentImages,
		GoogleAnalyticsCode: inp.GoogleAnalyticsCode,
		LogoURL:             inp.LogoURL,
		DisableRegistration: inp.DisableRegistration,
	}

	if inp.Pages != nil {
		updateInput.Pages = &domain.UpdateSchoolSettingsPages{
			Confidential:      inp.Pages.Confidential,
			ServiceAgreement:  inp.Pages.ServiceAgreement,
			NewsletterConsent: inp.Pages.NewsletterConsent,
		}
	}

	if inp.ContactInfo != nil {
		updateInput.ContactInfo = &domain.UpdateSchoolSettingsContactInfo{
			Email:              inp.ContactInfo.Email,
			RegistrationNumber: inp.ContactInfo.RegistrationNumber,
			Address:            inp.ContactInfo.Address,
			BusinessName:       inp.ContactInfo.BusinessName,
			Phone:              inp.ContactInfo.Phone,
		}
	}

	if err := h.services.Schools.UpdateSettings(c.Request.Context(), school.ID, updateInput); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

type connectFondyInput struct {
	MerchantID       string `json:"merchantId"`
	MerchantPassword string `json:"merchantPassword"`
}

// @Summary Admin Connect Fondy
// @Security AdminAuth
// @Tags admins-school
// @Description admin connect fondy
// @ModuleID adminConnectFondy
// @Accept  json
// @Produce  json
// @Param input body connectFondyInput true "update school settings"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/school/settings/fondy [put]
func (h *Handler) adminConnectFondy(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	var inp connectFondyInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	if err := h.services.Schools.ConnectFondy(c.Request.Context(), service.ConnectFondyInput{
		SchoolID:         school.ID,
		MerchantID:       inp.MerchantID,
		MerchantPassword: inp.MerchantPassword,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

type connectSendPulseInput struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	ListID string `json:"listId"`
}

// @Summary Admin Connect Fondy
// @Security AdminAuth
// @Tags admins-school
// @Description admin connect fondy
// @ModuleID adminConnectFondy
// @Accept  json
// @Produce  json
// @Param input body connectSendPulseInput true "update school settings"
// @Success 200 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/school/settings/sendpulse [put]
func (h *Handler) adminConnectSendPulse(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	var inp connectSendPulseInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	if err := h.services.Schools.ConnectSendPulse(c.Request.Context(), service.ConnectSendPulseInput{
		SchoolID: school.ID,
		ID:       inp.ID,
		Secret:   inp.Secret,
		ListID:   inp.ListID,
	}); err != nil {
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
// @Param skip query int false "skip"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param status query string false "status"
// @Param dateFrom query string false "dateFrom"
// @Param dateTo query string false "dateTo"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/orders [get]
func (h *Handler) adminGetOrders(c *gin.Context) {
	var query domain.GetOrdersQuery
	if err := c.Bind(&query); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	orders, count, err := h.services.Orders.GetBySchool(c.Request.Context(), school.ID, query)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data:  orders,
		Count: count,
	})
}

type orderStatusInput struct {
	Status string `json:"status" binding:"required"`
}

func (i orderStatusInput) validate() error {
	switch i.Status {
	case domain.OrderStatusPaid, domain.OrderStatusCanceled:
		return nil
	default:
		return errors.New("incorrect status")
	}
}

// @Summary Admin Update Order
// @Security AdminAuth
// @Tags admins-orders
// @Description admin update order status
// @ModuleID adminUpdateOrderStatus
// @Accept  json
// @Param id path string true "promocode id"
// @Param input body orderStatusInput true "update school settings"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/orders/{id} [put]
func (h *Handler) adminUpdateOrderStatus(c *gin.Context) {
	id, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	var inp orderStatusInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	if err := inp.validate(); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	if err := h.services.Orders.SetStatus(c.Request.Context(), id, inp.Status); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

func toPackagesResponse(pkgs []domain.Package) []packageResponse {
	out := make([]packageResponse, len(pkgs))

	for i, pkg := range pkgs {
		out[i] = packageResponse{
			ID:      pkg.ID,
			Name:    pkg.Name,
			Modules: toPackageModules(pkg.Modules),
		}
	}

	return out
}

func toPackageModules(modules []domain.Module) []packageModule {
	out := make([]packageModule, len(modules))

	for i, module := range modules {
		out[i] = packageModule{
			ID:   module.ID,
			Name: module.Name,
		}
	}

	return out
}

func toOffersReponse(offers []domain.Offer) []offerShortReponse {
	out := make([]offerShortReponse, len(offers))

	for i, offer := range offers {
		out[i] = offerShortReponse{
			ID:   offer.ID,
			Name: offer.Name,
		}
	}

	return out
}

func stringArrayToObjectId(stringIds []string) ([]primitive.ObjectID, error) {
	var err error

	ids := make([]primitive.ObjectID, len(stringIds))

	for i, id := range stringIds {
		ids[i], err = primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}
