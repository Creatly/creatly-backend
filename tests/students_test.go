package tests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/pkg/email"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
)

const (
	verificationCode = "CODE1234"
	schoolIdStr      = "6008090356fabf0d0b34635b"
)

func (s *APITestSuite) TestStudentSignUp() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))

	r := s.Require()

	name, studentEmail, password := "Test Student", "test@test.com", "qwerty123"
	signUpData := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, name, studentEmail, password)

	s.mocks.otpGenerator.On("RandomSecret", 8).Return(verificationCode)
	s.mocks.emailProvider.On("AddEmailToList", email.AddEmailInput{
		Email:  studentEmail,
		ListID: listId,
		Variables: map[string]string{
			"name":             name,
			"verificationCode": verificationCode,
		},
	}).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/students/sign-up", bytes.NewBuffer([]byte(signUpData)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	var student domain.Student
	err := s.db.Collection("students").FindOne(context.Background(), bson.M{"email": studentEmail}).Decode(&student)
	s.NoError(err)

	r.Equal(name, student.Name)
	r.Equal(s.hasher.Hash(password), student.Password)
	r.Equal(false, student.Verification.Verified)
	r.Equal(verificationCode, student.Verification.Code)
}

func (s *APITestSuite) TestStudentSignInNotVerified() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test2@test.com", "qwerty123"
	s.db.Collection("students").InsertOne(context.Background(), domain.Student{
		Email:    studentEmail,
		Password: s.hasher.Hash(password),
	})

	signUpData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, studentEmail, password)
	req, _ := http.NewRequest("POST", "/api/v1/students/sign-in", bytes.NewBuffer([]byte(signUpData)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestStudentVerify() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test3@test.com", "qwerty123"
	s.db.Collection("students").InsertOne(context.Background(), domain.Student{
		Email:        studentEmail,
		Password:     s.hasher.Hash(password),
		Verification: domain.Verification{Code: "CODE4321"},
	})

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/v1/students/verify/%s", "CODE4321"), nil)
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var student domain.Student
	err := s.db.Collection("students").FindOne(context.Background(), bson.M{"email": studentEmail}).Decode(&student)
	s.NoError(err)

	r.Equal(true, student.Verification.Verified)
	r.Equal("", student.Verification.Code)
}

func (s *APITestSuite) TestStudentSignInVerified() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	studentEmail, password := "test4@test.com", "qwerty123"
	schoolId, _ := primitive.ObjectIDFromHex(schoolIdStr)
	s.db.Collection("students").InsertOne(context.Background(), domain.Student{
		Email:        studentEmail,
		Password:     s.hasher.Hash(password),
		SchoolID:     schoolId,
		Verification: domain.Verification{Verified: true},
	})

	signUpData := fmt.Sprintf(`{"email":"%s","password":"%s"}`, studentEmail, password)

	req, _ := http.NewRequest("POST", "/api/v1/students/sign-in", bytes.NewBuffer([]byte(signUpData)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
}
