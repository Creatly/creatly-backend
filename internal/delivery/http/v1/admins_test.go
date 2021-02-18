package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/service"
	mock_service "github.com/zhashkevych/courses-backend/internal/service/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler_adminUpdateSchoolSettings(t *testing.T) {
	type mockBehavior func(r *mock_service.MockSchools, input service.UpdateSchoolSettingsInput)

	school := domain.School{
		ID: primitive.NewObjectID(),
		Settings: domain.Settings{
			Domain: "localhost",
		},
	}

	tests := []struct {
		name         string
		body         string
		school       domain.School
		input        service.UpdateSchoolSettingsInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			body:   fmt.Sprintf(`{"color": "black", "pages": {"confidential": "some confidential info"}}`),
			school: school,
			input: service.UpdateSchoolSettingsInput{
				SchoolID: school.ID,
				Color:    "black",
				Pages: &domain.Pages{
					Confidential: "some confidential info",
				},
			},
			mockBehavior: func(r *mock_service.MockSchools, input service.UpdateSchoolSettingsInput) {
				r.EXPECT().UpdateSettings(context.Background(), input).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:         "invalid input",
			body:         fmt.Sprintf(`{wrong}`),
			school:       school,
			mockBehavior: func(r *mock_service.MockSchools, input service.UpdateSchoolSettingsInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "service error",
			body:   fmt.Sprintf(`{"color": "black", "pages": {"confidential": "some confidential info"}}`),
			school: school,
			input: service.UpdateSchoolSettingsInput{
				SchoolID: school.ID,
				Color:    "black",
				Pages: &domain.Pages{
					Confidential: "some confidential info",
				},
			},
			mockBehavior: func(r *mock_service.MockSchools, input service.UpdateSchoolSettingsInput) {
				r.EXPECT().UpdateSettings(context.Background(), input).Return(errors.New("failed to update school settings"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to update school settings"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			s := mock_service.NewMockSchools(c)
			tt.mockBehavior(s, tt.input)

			services := &service.Services{Schools: s}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.PUT("/admins/school/settings", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminUpdateSchoolSettings)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/admins/school/settings",
				bytes.NewBufferString(tt.body))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}
