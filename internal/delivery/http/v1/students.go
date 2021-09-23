package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			authenticated.GET("/modules/:id/content", h.studentGetModuleContent)
			authenticated.GET("/modules/:id/offers", h.studentGetModuleOffers)
			authenticated.POST("/modules/:id/survey", h.studentSubmitSurvey)
			authenticated.POST("/lessons/:id/finished", h.studentSetLessonFinished)
			authenticated.POST("/orders", h.studentCreateOrder)
			authenticated.GET("/orders/:id/payment", h.studentGeneratePaymentLink)
			authenticated.GET("/account", h.studentGetAccount)
		}
	}
}

type studentSignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Verified bool   `json:"verified"`
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

	schoolDomain, err := getDomainFromContext(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Students.SignUp(c.Request.Context(), service.StudentSignUpInput{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     inp.Password,
		SchoolID:     school.ID,
		SchoolDomain: schoolDomain,
		Verified:     inp.Verified,
	}); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

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

	res, err := h.services.Students.SignIn(c.Request.Context(), service.SchoolSignInInput{
		SchoolID: school.ID,
		Email:    inp.Email,
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		if errors.Is(err, domain.ErrStudentBlocked) {
			newResponse(c, http.StatusForbidden, err.Error())

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
		if errors.Is(err, domain.ErrStudentBlocked) {
			newResponse(c, http.StatusForbidden, err.Error())

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
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, response{"success"})
}

// @Summary Student Get Content By Module ID
// @Security StudentsAuth
// @Tags students-courses
// @Description student get content by module id
// @ModuleID studentGetModuleContent
// @Accept  json
// @Produce  json
// @Param id path string true "module id"
// @Success 200 {object} domain.ModuleContent
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/modules/{id}/content [get]
func (h *Handler) studentGetModuleContent(c *gin.Context) {
	moduleId, err := parseIdFromPath(c, "id")
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

	content, err := h.services.Students.GetModuleContent(c.Request.Context(), school.ID, studentId, moduleId)
	if err != nil {
		if errors.Is(err, domain.ErrModuleIsNotAvailable) {
			newResponse(c, http.StatusForbidden, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, content)
}

type submitSurveyInput struct {
	Answers []surveyAnswer `json:"answers"`
}

type surveyAnswer struct {
	QuestionID string `json:"questionId"`
	Answer     string `json:"answer"`
}

// @Summary Student Submit Survey by Module ID
// @Security StudentsAuth
// @Tags students-courses
// @Description student submit survey by module id
// @ModuleID studentSetLessonFinished
// @Accept  json
// @Produce  json
// @Param input body submitSurveyInput true "survey answers"
// @Param id path string true "module id"
// @Success 200 {string} string "ok"
// @Failure 400,403 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/modules/{id}/survey [post]
func (h *Handler) studentSubmitSurvey(c *gin.Context) {
	moduleId, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	var inp submitSurveyInput
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

		return
	}

	answers, err := toSurveyAnswers(inp.Answers)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")

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

	if err := h.services.Surveys.SaveStudentAnswers(c.Request.Context(), service.SaveStudentAnswersInput{
		StudentID: studentId,
		SchoolID:  school.ID,
		ModuleID:  moduleId,
		Answers:   answers,
	}); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

// @Summary Student Set Lesson As Finished By LessonID
// @Security StudentsAuth
// @Tags students-courses
// @Description student set lesson as finished by lesson id
// @ModuleID studentSetLessonFinished
// @Accept  json
// @Produce  json
// @Param id path string true "lesson id"
// @Success 200 {string} string "ok"
// @Failure 400,403 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/lessons/{id}/finished [post]
func (h *Handler) studentSetLessonFinished(c *gin.Context) {
	lessonId, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	studentId, err := getStudentId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	if err := h.services.Students.SetLessonFinished(c.Request.Context(), studentId, lessonId); err != nil {
		if errors.Is(err, domain.ErrModuleIsNotAvailable) {
			newResponse(c, http.StatusForbidden, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}

type studentOffer struct {
	ID            primitive.ObjectID `json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Price         price              `json:"price"`
	Benefits      []string           `json:"benefits"`
	PaymentMethod struct {
		UsesProvider bool `json:"usesProvider"`
	} `json:"paymentMethod"`
}

type price struct {
	Value    uint   `json:"value" binding:"required,min=1"`
	Currency string `json:"currency" binding:"required,min=3"` // TODO validate currency input
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
		Benefits:    offer.Benefits,
		Price: price{
			Value:    offer.Price.Value,
			Currency: offer.Price.Currency,
		},
		PaymentMethod: struct {
			UsesProvider bool `json:"usesProvider"`
		}{
			offer.PaymentMethod.UsesProvider,
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
	moduleId, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

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

	c.JSON(http.StatusOK, dataResponse{Data: toStudentOffers(offers)})
}

type createOrderInput struct {
	OfferId string `json:"offerId" binding:"required"`
	PromoId string `json:"promoId"`
}

type createOrderResponse struct {
	OrderId string `json:"orderId"`
}

// @Summary Student CreateOrder
// @Security StudentsAuth
// @Tags students-orders
// @Description student create order
// @ModuleID studentCreateOrder
// @Accept  json
// @Produce  json
// @Param input body createOrderInput true "order info"
// @Success 200 {object} createOrderResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/orders [post]
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

	id, err := h.services.Orders.Create(c.Request.Context(), studentId, offerId, promoId)
	if err != nil {
		switch err {
		case domain.ErrPromoNotFound, domain.ErrOfferNotFound, domain.ErrUserNotFound, domain.ErrPromocodeExpired:
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		default:
			newResponse(c, http.StatusInternalServerError, err.Error())

			return
		}
	}

	c.JSON(http.StatusOK, createOrderResponse{id.Hex()})
}

type generatePaymentLinkResponse struct {
	URL string `json:"url"`
}

// @Summary Student Generate Payment Link
// @Security StudentsAuth
// @Tags students-orders
// @Description student generate order payment link
// @ModuleID studentGeneratePaymentLink
// @Accept  json
// @Produce  json
// @Param id path string true "order id"
// @Success 200 {object} generatePaymentLinkResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/orders/{id}/payment [get]
func (h *Handler) studentGeneratePaymentLink(c *gin.Context) {
	orderId, err := parseIdFromPath(c, "id")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	url, err := h.services.Payments.GeneratePaymentLink(c.Request.Context(), orderId)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentProviderNotUsed) {
			newResponse(c, http.StatusBadRequest, err.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, generatePaymentLinkResponse{url})
}

type studentAccountResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// @Summary Student Get Account Info
// @Security StudentsAuth
// @Tags students-account
// @Description student get account info
// @ModuleID studentGetAccount
// @Accept  json
// @Produce  json
// @Success 200 {object} studentAccountResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /students/account [get]
func (h *Handler) studentGetAccount(c *gin.Context) {
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

	student, err := h.services.Students.GetById(c.Request.Context(), school.ID, studentId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.JSON(http.StatusOK, studentAccountResponse{
		Name:  student.Name,
		Email: student.Email,
	})
}

func toSurveyAnswers(answers []surveyAnswer) ([]domain.SurveyAnswer, error) {
	res := make([]domain.SurveyAnswer, len(answers))

	for i := range answers {
		id, err := primitive.ObjectIDFromHex(answers[i].QuestionID)
		if err != nil {
			return nil, err
		}

		res[i] = domain.SurveyAnswer{
			QuestionID: id,
			Answer:     answers[i].Answer,
		}
	}

	return res, nil
}
