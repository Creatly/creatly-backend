package tests

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func (s *APITestSuite) TestAdminGetAllCourses() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	req, _ := http.NewRequest("GET", "/api/v1/admins/courses", nil)
	req.Header.Set("Content-type", "application/json")

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

	r.Equal(1, len(respCourses.Data))
}
