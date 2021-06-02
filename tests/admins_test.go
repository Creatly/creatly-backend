package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *APITestSuite) TestAdminGetAllCourses() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	id := primitive.NewObjectID()
	studentEmail, password := "test4@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	_, err = s.db.Collection("admins").InsertOne(context.Background(), domain.Student{
		ID:       id,
		Email:    studentEmail,
		Password: passwordHash,
		SchoolID: school.ID,
	})
	s.NoError(err)

	jwt, err := s.getJwt(id)
	s.NoError(err)

	req, _ := http.NewRequest("GET", "/api/v1/admins/courses", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var respCourses struct {
		Data []courseResponse `json:"data"`
	}

	respData, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	err = json.Unmarshal(respData, &respCourses)
	s.NoError(err)

	r.Equal(2, len(respCourses.Data))
}

func (s *APITestSuite) TestAdminCreateCourse() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// populate DB data
	id := primitive.NewObjectID()
	studentEmail, password := "test4@test.com", "qwerty123"
	passwordHash, err := s.hasher.Hash(password)
	s.NoError(err)

	_, err = s.db.Collection("admins").InsertOne(context.Background(), domain.Student{
		ID:       id,
		Email:    studentEmail,
		Password: passwordHash,
		SchoolID: school.ID,
	})
	s.NoError(err)

	jwt, err := s.getJwt(id)
	s.NoError(err)

	adminCourseName := "admin course test name"

	name := fmt.Sprintf(`{"name":"%s"}`, adminCourseName)

	req, _ := http.NewRequest("POST", "/api/v1/admins/courses", strings.NewReader(name))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)
}
