package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) initCoursesRoutes(api *gin.RouterGroup) {
	courses := api.Group("/courses", h.setSchoolFromRequest)
	{
		courses.GET("", h.getAllCourses)
		courses.GET("/:id", h.getCourseById)
		courses.GET("/:id/offers", h.getCourseOffers)
	}
}

// @Summary  Get All Courses
// @Tags courses
// @Description  get all courses
// @ModuleID getAllCourses
// @Accept  json
// @Produce  json
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /courses [get]
func (h *Handler) getAllCourses(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	// Return only published courses
	courses := make([]domain.Course, 0)

	for _, course := range school.Courses {
		if course.Published {
			courses = append(courses, course)
		}
	}

	c.JSON(http.StatusOK, dataResponse{Data: courses})
}

type getCourseByIdResponse struct {
	Course  domain.Course `json:"course"`
	Modules []module      `json:"modules"`
}

type module struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Position uint               `json:"position" bson:"position"`
	Lessons  []lesson           `json:"lessons" bson:"lessons"`
}

type lesson struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Position uint               `json:"position" bson:"position"`
}

func newGetCourseByIdResponse(course domain.Course, courseModules []domain.Module) getCourseByIdResponse {
	modules := make([]module, len(courseModules))

	for i := range courseModules {
		modules[i].ID = courseModules[i].ID
		modules[i].Name = courseModules[i].Name
		modules[i].Position = courseModules[i].Position
		modules[i].Lessons = toLessons(courseModules[i].Lessons)
	}

	return getCourseByIdResponse{
		Course:  course,
		Modules: modules,
	}
}

func toLessons(lessons []domain.Lesson) []lesson {
	out := make([]lesson, 0)

	for _, l := range lessons {
		if l.Published {
			out = append(out, lesson{
				ID:       l.ID,
				Name:     l.Name,
				Position: l.Position,
			})
		}
	}

	return out
}

// @Summary Get Course By ModuleID
// @Tags courses
// @Description  get course by id
// @ModuleID getCourseById
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} domain.Course
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /courses/{id} [get]
func (h *Handler) getCourseById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	course, err := studentGetSchoolCourse(school, id)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	modules, err := h.services.Modules.GetPublishedByCourseId(c.Request.Context(), course.ID)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, newGetCourseByIdResponse(course, modules))
}

func studentGetSchoolCourse(school domain.School, courseId string) (domain.Course, error) {
	var searchedCourse domain.Course

	for _, course := range school.Courses {
		if course.Published && course.ID.Hex() == courseId {
			searchedCourse = course
		}
	}

	if searchedCourse.ID.IsZero() {
		return domain.Course{}, errors.New("not found")
	}

	return searchedCourse, nil
}

// @Summary Get Course Offers
// @Tags courses
// @Description  get course offers
// @ModuleID getCourseOffers
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /courses/{id}/offers [get]
func (h *Handler) getCourseOffers(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "empty id param")

		return
	}

	courseId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid id param")

		return
	}

	offers, err := h.services.Offers.GetByCourse(c.Request.Context(), courseId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, dataResponse{Data: offers})
}
