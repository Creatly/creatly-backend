package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// TODO: review response error messages

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
				offers.POST("/", h.adminCreateOffer)
				offers.GET("/", h.adminGetAllOffers)
				offers.GET("/:id", h.adminGetOfferById)
				offers.PUT("/:id", h.adminUpdateOffer)
				offers.DELETE("/:id", h.adminDeleteOffer)
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
// @ModuleID adminCreateCourse
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
// @ModuleID adminGetAllCourses
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

type updateCourseInput struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
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
		"id": moduleId, // TODO create structure for id responses
	})
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

	err = h.modulesService.Delete(c.Request.Context(), id)
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
// @Success 200 {string} string "ok"
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

type createLessonInput struct {
	Name     string `json:"name" binding:"required,min=5"`
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
// @Success 201 {string} string "id"
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

	lessonId, err := h.modulesService.AddLesson(c.Request.Context(), service.AddLessonInput{
		ModuleID: id,
		Name:     inp.Name,
		Position: inp.Position,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": lessonId, // TODO create structure for id responses
	})
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

	lesson, err := h.modulesService.GetLesson(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, lesson)
}

type updateLessonInput struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	Position  *uint   `json:"position"`
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

	err := h.modulesService.UpdateLesson(c.Request.Context(), service.UpdateLessonInput{
		LessonID:  id,
		Name:      inp.Name,
		Content:   inp.Content,
		Position:  inp.Position,
		Published: inp.Published,
	})
	if err != nil {
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

	err = h.modulesService.DeleteLesson(c.Request.Context(), id)
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
// @Success 201 {string} string "id"
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

	moduleId, err := h.packagesService.Create(c.Request.Context(), service.CreatePackageInput{
		CourseID:    id,
		Name:        inp.Name,
		Description: inp.Description,
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": moduleId, // TODO create structure for id responses
	})
}

type getAllPackagesResponse struct {
	Packages []domain.Package `json:"packages"`
}

// @Summary Admin Get All Course Packages
// @Security AdminAuth
// @Tags admins-packages
// @Description admin get all course packages
// @ModuleID adminGetAllPackages
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} getAllPackagesResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/courses/{id}/packages [get]
func (h *Handler) adminGetAllPackages(c *gin.Context) {
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

	packages, err := h.packagesService.GetByCourse(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	c.JSON(http.StatusOK, getAllPackagesResponse{Packages: packages})
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

	pkg, err := h.packagesService.GetById(c.Request.Context(), id)
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

	if err := h.packagesService.Update(c.Request.Context(), service.UpdatePackageInput{
		ID:          id,
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

	err = h.packagesService.Delete(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "invalid id param")
		return
	}

	c.Status(http.StatusOK)
}

type createOfferInput struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description"`
	Price       price  `json:"price" binding:"required"`
}

// @Summary Admin Create Offer
// @Security AdminAuth
// @Tags admins-offers
// @Description admin create offer
// @ModuleID adminCreateOffer
// @Accept  json
// @Produce  json
// @Param input body createOfferInput true "package info"
// @Success 201 {string} string "id"
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

	id, err := h.offersService.Create(c.Request.Context(), service.CreateOfferInput{
		SchoolID:    school.ID,
		Name:        inp.Name,
		Description: inp.Description,
		Price: domain.Price{
			Value:    inp.Price.Value,
			Currency: inp.Price.Currency,
		},
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id, // TODO create structure for id responses
	})
}

type getAllOffersResponse struct {
	Offers []domain.Offer `json:"offers"`
}

// @Summary Admin Get All Offers
// @Security AdminAuth
// @Tags admins-offers
// @Description admin get all offers
// @ModuleID adminGetAllOffers
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllOffersResponse
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

	offers, err := h.offersService.GetAll(c.Request.Context(), school.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllOffersResponse{offers})
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

	offer, err := h.offersService.GetById(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, offer)
}

type updateOfferInput struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
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

	updateInput := service.UpdateOfferInput{
		ID:          id,
		Name:        inp.Name,
		Description: inp.Description,
		Packages:    inp.Packages,
	}

	if inp.Price != nil {
		updateInput.Price = &domain.Price{
			Value:    inp.Price.Value,
			Currency: inp.Price.Currency,
		}
	}

	if err := h.offersService.Update(c.Request.Context(), updateInput); err != nil {
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

	err = h.offersService.Delete(c.Request.Context(), id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
