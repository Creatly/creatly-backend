package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *Handler) initStudentsRoutes(api *gin.RouterGroup) {
	students := api.Group("/students", h.setSchoolFromRequest)
	{
		students.POST("/sign-up", h.studentSignUp)
		students.POST("/sign-in", h.studentSignIn)
		students.POST("/auth/refresh", h.studentRefresh)
		students.POST("/verify/:hash", h.studentVerify)
		students.GET("/courses", h.studentGetAllCourses)
		students.GET("/courses/:id", h.studentGetCourseById)

		_ = students.Group("/", h.userIdentity)
		{

		}
	}
}

type studentSignUpInput struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RegisterSource string `json:"registerSource"`
}

// @Summary Student SignUp
// @Tags students
// @Description create student account
// @ID studentSignUp
// @Accept  json
// @Produce  json
// @Param input body studentSignUpInput true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /students/sign-up [post]
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

// @Summary Student SignIn
// @Tags students
// @Description student sign in
// @ID studentSignIn
// @Accept  json
// @Produce  json
// @Param input body studentSignInInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /students/sign-in [post]
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

// @Summary Student Refresh Tokens
// @Tags students
// @Description student refresh tokens
// @Accept  json
// @Produce  json
// @Param input body refreshInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /students/refresh [post]
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

// @Summary Student Verify Registration
// @Tags students
// @Description student verify registration
// @ID studentVerify
// @Accept  json
// @Produce  json
// @Param code path string true "verification code"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /students/verify/{code} [post]
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

// @Summary Student Get All Courses
// @Tags students
// @Description student get all courses
// @ID studentGetAllCourses
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Course
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /students/courses [get]
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

type getCourseByIdResponse struct {
	Course  domain.Course `json:"course"`
	Modules []module      `json:"modules"`
}

type module struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Position int                `json:"position" bson:"position"`
	Lessons  []lesson           `json:"lessons" bson:"lessons"`
}

type lesson struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Position int                `json:"position" bson:"position"`
}

func newGetCourseByIdResponse(course domain.Course, courseModules []domain.Module) getCourseByIdResponse {
	modules := make([]module, len(courseModules))

	for i := range modules {
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

// @Summary Student Get Course By ID
// @Tags students
// @Description student get course by id
// @ID studentGetCourseById
// @Accept  json
// @Produce  json
// @Param id path string true "course id"
// @Success 200 {object} domain.Course
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /students/courses/{id} [get]
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

	var searchedCourse domain.Course
	for _, course := range school.Courses {
		if course.Published && course.ID.Hex() == id {
			searchedCourse = course
		}
	}

	if searchedCourse.ID.IsZero() {
		newErrorResponse(c, http.StatusBadRequest, "not found")
		return
	}

	modules, err := h.coursesService.GetCourseModules(c.Request.Context(), searchedCourse.ID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newGetCourseByIdResponse(searchedCourse, modules))
}
