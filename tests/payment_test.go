package tests

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/pkg/email"
	"github.com/zhashkevych/creatly-backend/pkg/payment/fondy"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *APITestSuite) TestFondyCallbackApproved() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentId := primitive.NewObjectID()
	studentEmail := "payment@test.com"
	studentName := "Test Payment"
	offerName := "Test Offer"
	_, err := s.db.Collection("students").InsertOne(context.Background(), domain.Student{
		ID:           studentId,
		Email:        studentEmail,
		Name:         studentName,
		SchoolID:     school.ID,
		Verification: domain.Verification{Verified: true},
	})
	s.NoError(err)

	id, _ := primitive.ObjectIDFromHex("6008153f3dab4fb0573d1f96")
	_, err = s.db.Collection("orders").InsertOne(context.Background(), domain.Order{
		ID:       id,
		SchoolID: school.ID,
		Offer:    domain.OrderOfferInfo{ID: offers[0].(domain.Offer).ID, Name: offerName},
		Student:  domain.StudentInfoShort{ID: studentId, Email: studentEmail, Name: studentName},
		Status:   "created",
	})
	s.NoError(err)

	s.mocks.emailSender.On("Send", email.SendEmailInput{
		To:      studentEmail,
		Subject: "Покупка прошла успешно!",
		Body: fmt.Sprintf(`<h1>%s, спасибо большое за покупку "%s"!</h1>
<br>
<p>Надеюсь данный материал будет тебе полезен и интересен!</p>
<p>Если у тебя возникают вопросы или ты хочешь поделиться своим отзывом - пиши мне письмо на <a href="mailto:zhashkevychmaksim@gmail.com">zhashkevychmaksim@gmail.com</a>.</p>
<p>Мне крайне важен твой отзыв, чтобы улучшать материалы и делать курс максимально полезным!</p>

<br><br>

<p><i>С уважением, Максим</i></p>`, studentName, offerName),
	}).Return(nil)

	file, err := ioutil.ReadFile("./fixtures/callback_approved.json")
	s.NoError(err)

	req, _ := http.NewRequest("POST", "/api/v1/callback/fondy", bytes.NewBuffer(file))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", fondy.UserAgent)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// Get Paid Lessons After Callback
	r = s.Require()

	jwt, err := s.getJwt(studentId)
	s.NoError(err)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/students/modules/%s/content", modules[1].(domain.Module).ID.Hex()), nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
}

func (s *APITestSuite) TestFondyCallbackDeclined() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentId := primitive.NewObjectID()
	_, err := s.db.Collection("students").InsertOne(context.Background(), domain.Student{
		ID:           studentId,
		SchoolID:     school.ID,
		Verification: domain.Verification{Verified: true},
	})
	s.NoError(err)

	id, _ := primitive.ObjectIDFromHex("6008153f3dab4fb0573d1f97")
	_, err = s.db.Collection("orders").InsertOne(context.Background(), domain.Order{
		ID:       id,
		SchoolID: school.ID,
		Offer:    domain.OrderOfferInfo{ID: offers[0].(domain.Offer).ID},
		Student:  domain.StudentInfoShort{ID: studentId},
		Status:   "created",
	})
	s.NoError(err)

	file, err := ioutil.ReadFile("./fixtures/callback_declined.json")
	s.NoError(err)

	req, _ := http.NewRequest("POST", "/api/v1/callback/fondy", bytes.NewBuffer(file))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", fondy.UserAgent)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// Get Paid Lessons After Callback
	r = s.Require()

	jwt, err := s.getJwt(studentId)
	s.NoError(err)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/students/modules/%s/content", modules[1].(domain.Module).ID.Hex()), nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusForbidden, resp.Result().StatusCode)
}
