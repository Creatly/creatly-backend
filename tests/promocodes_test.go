package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/creatly-backend/internal/domain"
)

func (s *APITestSuite) TestGetPromoCode() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/promocodes/%s", promocodes[0].(domain.PromoCode).Code), nil)
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
}

func (s *APITestSuite) TestGetPromoCodeInvalid() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	req, _ := http.NewRequest("GET", "/api/v1/promocodes/CODE123", nil)
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}
