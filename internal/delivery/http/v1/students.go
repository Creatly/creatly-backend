package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

// TODO: return time.Time in RFC3339

func (h *Handler) initStudentsRoutes(api *gin.RouterGroup) {
	students := api.Group("/students", h.setSchoolFromRequest)
	{
		students.POST("/sign-up", h.studentSignUp)
		students.POST("/sign-in", h.studentSignIn)
		students.POST("/auth/refresh", h.studentRefresh)
		students.POST("/verify/:code", h.studentVerify)
		students.GET("/courses", h.getAllCourses)
		students.GET("/courses/:id", h.getCourseById)

		authenticated := students.Group("/", h.studentIdentity)
		{
			authenticated.GET("/modules/:id/lessons", h.studentGetModuleLessons)
			authenticated.GET("/modules/:id/offers", h.studentGetModuleOffers)
			authenticated.GET("/promocodes/:code", h.studentGetPromo)
			authenticated.POST("/order", h.studentCreateOrder)
		}
	}
}

type studentSignUpInput struct {
	Name           string `json:"name" binding:"required,min=2,max=64"`
	Email          string `json:"email" binding:"required,email,max=64"`
	Password       string `json:"password" binding:"required,min=8,max=64"`
	RegisterSource string `json:"registerSource" binding:"required,max=64"`
}

// @Summary Student SignUp
// @Tags students-auth
// @Description create student account
// @ID studentSignUp
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

	if err := h.studentsService.SignUp(c.Request.Context(), service.StudentSignUpInput{
		Name:           inp.Name,
		Email:          inp.Email,
		Password:       inp.Password,
		RegisterSource: inp.RegisterSource,
		SchoolID:       school.ID,
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
// @ID studentSignIn
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

	res, err := h.studentsService.SignIn(c.Request.Context(), service.SignInInput{
		SchoolID: school.ID,
		Email:    inp.Email,
		Password: inp.Password,
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

	res, err := h.studentsService.RefreshTokens(c.Request.Context(), school.ID, inp.Token)
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
// @ID studentVerify
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

	if err := h.studentsService.Verify(c.Request.Context(), code); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newResponse(c, http.StatusOK, "success")
}

type getModuleLessonsResponse struct {
	Lessons []domain.Lesson `json:"lessons"`
}

// @Summary Student Get Lessons By Module ID
// @Security StudentsAuth
// @Tags students-courses
// @Description student get lessons by module id
// @ID studentGetModuleLessons
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Success 200 {object} getModuleLessonsResponse
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

	lessons, err := h.studentsService.GetModuleLessons(c.Request.Context(), school.ID, studentId, moduleId)
	if err != nil {
		if err == service.ErrModuleIsNotAvailable {
			newResponse(c, http.StatusForbidden, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getModuleLessonsResponse{
		Lessons: lessons,
	})
}

type studentGetModuleOffersResponse struct {
	Offers []studentOffer `json:"offers"`
}

type studentOffer struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   string             `json:"createdAt" bson:"createdAt"`
	Price       price              `json:"price" bson:"price"`
}

type price struct {
	Value    int    `json:"value"`
	Currency string `json:"currency"`
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
		CreatedAt:   offer.CreatedAt.Format(time.RFC3339),
		Price: price{
			Value:    offer.Price.Value,
			Currency: offer.Price.Currency,
		},
	}
}

// @Summary Student Get Offers By Module ID
// @Security StudentsAuth
// @Tags students-courses
// @Description student get offers by module id
// @ID studentGetModuleOffers
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Success 200 {string} string "ok"
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

	offers, err := h.offersService.GetByModule(c.Request.Context(), school.ID, moduleId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, studentGetModuleOffersResponse{
		Offers: toStudentOffers(offers),
	})
}

// @Summary Student Get PromoCode By Code
// @Security StudentsAuth
// @Tags students-courses
// @Description student get promocode by code
// @ID studentGetPromo
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

	promocode, err := h.promoCodesService.GetByCode(c.Request.Context(), school.ID, code)
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
// @ID studentCreateOrder
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

	paymentLink, err := h.ordersService.Create(c.Request.Context(), studentId, offerId, promoId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, createOrderResponse{paymentLink})
}
