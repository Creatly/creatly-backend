package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/service"
	mock_service "github.com/zhashkevych/creatly-backend/internal/service/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler_getPromocode(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPromoCodes, schoolId primitive.ObjectID, code string, promocode domain.PromoCode)

	schoolId := primitive.NewObjectID()

	promocode := domain.PromoCode{
		Code:               "GOGOGO25",
		DiscountPercentage: 25,
	}

	setResponseBody := func(promocode domain.PromoCode) string {
		body, _ := json.Marshal(promocode)

		return string(body)
	}

	tests := []struct {
		name         string
		code         string
		schoolId     primitive.ObjectID
		promocode    domain.PromoCode
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:      "ok",
			code:      "GOGOGO25",
			schoolId:  schoolId,
			promocode: promocode,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId primitive.ObjectID, code string, promocode domain.PromoCode) {
				r.EXPECT().GetByCode(context.Background(), schoolId, code).Return(promocode, nil)
			},
			statusCode:   200,
			responseBody: setResponseBody(promocode),
		},
		{
			name:      "empty code",
			code:      "",
			schoolId:  schoolId,
			promocode: promocode,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId primitive.ObjectID, code string, promocode domain.PromoCode) {
			},
			statusCode:   404,
			responseBody: `404 page not found`,
		},
		{
			name:      "service error",
			code:      "GOGOGO25",
			schoolId:  schoolId,
			promocode: promocode,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId primitive.ObjectID, code string, promocode domain.PromoCode) {
				r.EXPECT().GetByCode(context.Background(), schoolId, code).Return(promocode, errors.New("failed to get promocode"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to get promocode"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			s := mock_service.NewMockPromoCodes(c)
			tt.mockBehavior(s, tt.schoolId, tt.code, tt.promocode)

			services := &service.Services{PromoCodes: s}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/promocodes/:code", func(c *gin.Context) {
				c.Set(schoolCtx, domain.School{
					ID: schoolId,
				})
			}, handler.getPromo)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/promocodes/%s", tt.code), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}
