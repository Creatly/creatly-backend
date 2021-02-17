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
		students.POST("/verify/:code", h.studentVerify)

		authenticated := students.Group("/", h.studentIdentity)
		{
			authenticated.GET("/courses", h.studentGetCourses)
			authenticated.GET("/modules/:id/lessons", h.studentGetModuleLessons)
			authenticated.GET("/modules/:id/offers", h.studentGetModuleOffers)
			authenticated.GET("/promocodes/:code", h.studentGetPromo)
			authenticated.POST("/order", h.studentCreateOrder)
		}
	}
}

type studentSignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

// @Summary Student SignUp
// @Tags students-auth
// @Description create student account
// @ModuleID studentSignUp
// @Accept  json
// @Produce  json
// @Param input body studentSignUpInput true "sign up info"
// @Success 201 {string} string "ok"
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/sign-up [post]
func (h *Handler) studentSignUp(c *gin.Context) {
	var inp studentSignUpInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.services.Students.SignUp(c.Request.Context(), service.StudentSignUpInput{
		Name:     inp.Name,
		Email:    inp.Email,
		Password: inp.Password,
		SchoolID: school.ID,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

type signInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// @Summary Student SignIn
// @Tags students-auth
// @Description student sign in
// @ModuleID studentSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/sign-in [post]
func (h *Handler) studentSignIn(c *gin.Context) {
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

	res, err := h.services.Students.SignIn(c.Request.Context(), service.SignInInput{
		SchoolID: school.ID,
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		if err == service.ErrUserNotFound {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
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
// @Tags students-auth
// @Description student refresh tokens
// @Accept  json
// @Produce  json
// @Param input body refreshInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/auth/refresh [post]
func (h *Handler) studentRefresh(c *gin.Context) {
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

	res, err := h.services.Students.RefreshTokens(c.Request.Context(), school.ID, inp.Token)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// @Summary Student Verify Registration
// @Tags students-auth
// @Description student verify registration
// @ModuleID studentVerify
// @Accept  json
// @Produce  json
// @Param code path string true "verification code"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/verify/{code} [post]
func (h *Handler) studentVerify(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newResponse(c, http.StatusBadRequest, "code is empty")
		return
	}

	if err := h.services.Students.Verify(c.Request.Context(), code); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(c, http.StatusOK, "success")
}

// @Summary Student Get Lessons By Module ModuleID
// @Security StudentsAuth
// @Tags students-courses
// @Description student get lessons by module id
// @ModuleID studentGetModuleLessons
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/modules/{id}/lessons [get]
func (h *Handler) studentGetModuleLessons(c *gin.Context) {
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

	studentId, err := getStudentId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	lessons, err := h.services.Students.GetModuleLessons(c.Request.Context(), school.ID, studentId, moduleId)
	if err != nil {
		if err == service.ErrModuleIsNotAvailable {
			newResponse(c, http.StatusForbidden, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataResponse{lessons})
}

type studentOffer struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       price              `json:"price" bson:"price"`
}

type price struct {
	Value    int    `json:"value" binding:"required,min=1"`
	Currency string `json:"currency" binding:"required,min=3"`
}

func toStudentOffers(offers []domain.Offer) []studentOffer {
	out := make([]studentOffer, len(offers))

	for i := range offers {
		out[i] = toStudentOffer(offers[i])
	}

	return out
}

func toStudentOffer(offer domain.Offer) studentOffer {
	return studentOffer{
		ID:          offer.ID,
		Name:        offer.Name,
		Description: offer.Description,
		Price: price{
			Value:    offer.Price.Value,
			Currency: offer.Price.Currency,
		},
	}
}

// @Summary Student Get Offers By Module ModuleID
// @Security StudentsAuth
// @Tags students-courses
// @Description student get offers by module id
// @ModuleID studentGetModuleOffers
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/modules/{id}/offers [get]
func (h *Handler) studentGetModuleOffers(c *gin.Context) {
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

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	offers, err := h.services.Offers.GetByModule(c.Request.Context(), school.ID, moduleId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataResponse{toStudentOffers(offers)})
}

// @Summary Student Get PromoCode By Code
// @Security StudentsAuth
// @Tags students-courses
// @Description student get promocode by code
// @ModuleID studentGetPromo
// @Accept  json
// @Produce  json
// @Param code path string true "code"
// @Success 200 {object} domain.PromoCode
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/promocodes/{code} [get]
func (h *Handler) studentGetPromo(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		newResponse(c, http.StatusBadRequest, "empty code param")
		return
	}

	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	promocode, err := h.services.PromoCodes.GetByCode(c.Request.Context(), school.ID, code)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, promocode)
}

type createOrderInput struct {
	OfferId string `json:"offerId" binding:"required"`
	PromoId string `json:"promoId"`
}

type createOrderResponse struct {
	PaymentLink string `json:"paymentLink"`
}

// @Summary Student CreateOrder
// @Security StudentsAuth
// @Tags students-courses
// @Description student create order
// @ModuleID studentCreateOrder
// @Accept  json
// @Produce  json
// @Param input body createOrderInput true "order info"
// @Success 200 {object} createOrderResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/order [post]
func (h *Handler) studentCreateOrder(c *gin.Context) {
	var inp createOrderInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	offerId, err := primitive.ObjectIDFromHex(inp.OfferId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid offer id")
		return
	}

	var promoId primitive.ObjectID

	if inp.PromoId != "" {
		var err error
		promoId, err = primitive.ObjectIDFromHex(inp.PromoId)
		if err != nil {
			newResponse(c, http.StatusBadRequest, "invalid promo id")
			return
		}
	}

	studentId, err := getStudentId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	paymentLink, err := h.services.Orders.Create(c.Request.Context(), studentId, offerId, promoId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, createOrderResponse{paymentLink})
}

// @Summary Student Get Opened Courses
// @Security StudentsAuth
// @Tags students-courses
// @Description student get opened courses
// @ModuleID studentGetCourses
// @Accept  json
// @Produce  json
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/courses/ [get]
func (h *Handler) studentGetCourses(c *gin.Context) {
	school, err := getSchoolFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	studentId, err := getStudentId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	courses, err := h.services.Students.GetAvailableCourses(c.Request.Context(), school, studentId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dataResponse{courses})
}
