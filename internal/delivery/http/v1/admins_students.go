package v1

import (
	"net/http"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type studentResponse struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Email        string             `json:"email"`
	RegisteredAt time.Time          `json:"registeredAt"`
	LastVisitAt  time.Time          `json:"lastVisitAt"`
	Verified     bool               `json:"verified"`
}

func toStudentsResponse(students []domain.Student) []studentResponse {
	out := make([]studentResponse, len(students))
	for i, student := range students {
		out[i].ID = student.ID
		out[i].Name = student.Name
		out[i].Email = student.Email
		out[i].RegisteredAt = student.RegisteredAt
		out[i].LastVisitAt = student.LastVisitAt
		out[i].Verified = student.Verification.Verified
	}

	return out
}

// @Summary Admin Get Students
// @Security AdminAuth
// @Tags admins-students
// @Description admin get all students
// @ModuleID adminGetStudents
// @Accept  json
// @Produce  json
// @Param skip query int false "skip"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param verified query bool false "verified"
// @Param registerDateFrom query string false "registerDateFrom"
// @Param registerDateTo query string false "registerDateTo"
// @Param lastVisitDateFrom query string false "lastVisitDateFrom"
// @Param lastVisitDateTo query string false "registerDateTo"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/students [get]
func (h *Handler) adminGetStudents(c *gin.Context) {
	var query domain.GetStudentsQuery
	if err := c.Bind(&query); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	students, count, err := h.services.Students.GetBySchool(c.Request.Context(), school.ID, query)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{
		Data:  toStudentsResponse(students),
		Count: count,
	})
}

// @Summary Admin Get Student By ID
// @Security AdminAuth
// @Tags admins-students
// @Description admin get student by id
// @ModuleID adminGetStudents
// @Accept  json
// @Produce  json
// @Param id path string true "student id"
// @Success 200 {object} domain.Student
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/students/{id} [get]
func (h *Handler) adminGetStudentById(c *gin.Context) {
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

	student, err := h.services.Students.GetById(c.Request.Context(), school.ID, id)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, student)
}

type createStudentInput struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// @Summary Admin Create Student
// @Security AdminAuth
// @Tags admins-students
// @Description admin create student
// @ModuleID adminCreateStudent
// @Accept  json
// @Produce  json
// @Param input body createStudentInput true "student info"
// @Success 200 {object} domain.Student
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/students [post]
func (h *Handler) adminCreateStudent(c *gin.Context) {
	var inp createStudentInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Admins.CreateStudent(c.Request.Context(), service.CreateStudentInput{
		SchoolID: school.ID,
		Name:     inp.Name,
		Email:    inp.Email,
		Password: inp.Password,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

type manageOfferPermissionInput struct {
	Available bool `json:"available"`
}

// @Summary Admin Give Student Access to Offer
// @Security AdminAuth
// @Tags admins-students
// @Description admin give student access to offer
// @ModuleID adminManageOfferPermission
// @Accept  json
// @Produce  json
// @Param input body manageOfferPermissionInput true "permission type"
// @Param id path string true "student id"
// @Param offerId path string true "offer id"
// @Success 200 {object} domain.Student
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/students/{id}/offers/{offerId} [patch]
func (h *Handler) adminManageOfferPermission(c *gin.Context) {
	var inp manageOfferPermissionInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	studentId, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	offerId, err := parseIdFromPath(c, "offerId")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	offer, err := h.services.Offers.GetById(c.Request.Context(), offerId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if inp.Available {
		if err := h.services.Students.GiveAccessToOffer(c.Request.Context(), studentId, offer); err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())

			return
		}

		c.Status(http.StatusOK)

		return
	}

	if err := h.services.Students.RemoveAccessToOffer(c.Request.Context(), studentId, offer); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}
