package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

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
// @Tags admins-students
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

// @Summary Admin Give Student Access to Offer
// @Security AdminAuth
// @Tags admins-students
// @Description admin give student access to offer
// @ModuleID adminGiveAccessToOffer
// @Accept  json
// @Produce  json
// @Param id path string true "student id"
// @Param offerId path string true "offer id"
// @Success 200 {object} domain.Student
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/students/{id}/offers/{offerId} [post]
func (h *Handler) adminGiveAccessToOffer(c *gin.Context) {

}