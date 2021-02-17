package tests

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/pkg/email"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestSignUp() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	signUpData := `{"name":"Test User","email":"test@test.com","password":"qwerty123"}`
	s.mocks.otpGenerator.On("RandomSecret", 8).Return("CODE1234")
	s.mocks.emailProvider.On("AddEmailToList", email.AddEmailInput{
		Email:  "test@test.com",
		ListID: listId,
		Variables: map[string]string{
			"name":             "Test User",
			"verificationCode": "CODE1234",
		},
	}).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/students/sign-up", bytes.NewBuffer([]byte(signUpData)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)
}
