package v1

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/service"
	mock_service "github.com/zhashkevych/creatly-backend/internal/service/mocks"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler_adminUpdateSchoolSettings(t *testing.T) {
	type mockBehavior func(r *mock_service.MockSchools, schoolID primitive.ObjectID, input domain.UpdateSchoolSettingsInput)

	school := domain.School{
		ID: primitive.NewObjectID(),
		Settings: domain.Settings{
			Domains: []string{"localhost"},
		},
	}

	tests := []struct {
		name         string
		body         string
		school       domain.School
		input        domain.UpdateSchoolSettingsInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			body:   `{"color": "black", "pages": {"confidential": "some confidential info"}}`,
			school: school,
			input: domain.UpdateSchoolSettingsInput{
				Color: stringPtr("black"),
				Pages: &domain.UpdateSchoolSettingsPages{
					Confidential: stringPtr("some confidential info"),
				},
			},
			mockBehavior: func(r *mock_service.MockSchools, schoolID primitive.ObjectID, input domain.UpdateSchoolSettingsInput) {
				r.EXPECT().UpdateSettings(context.Background(), schoolID, input).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:   "invalid input",
			body:   `{wrong}`,
			school: school,
			mockBehavior: func(r *mock_service.MockSchools, schoolID primitive.ObjectID, input domain.UpdateSchoolSettingsInput) {
			},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "service error",
			body:   `{"color": "black", "pages": {"confidential": "some confidential info"}}`,
			school: school,
			input: domain.UpdateSchoolSettingsInput{
				Color: stringPtr("black"),
				Pages: &domain.UpdateSchoolSettingsPages{
					Confidential: stringPtr("some confidential info"),
				},
			},
			mockBehavior: func(r *mock_service.MockSchools, schoolID primitive.ObjectID, input domain.UpdateSchoolSettingsInput) {
				r.EXPECT().UpdateSettings(context.Background(), schoolID, input).Return(errors.New("failed to update school settings"))
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
			tt.mockBehavior(s, school.ID, tt.input)

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

func TestHandler_adminCreatePromocode(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPromoCodes, input service.CreatePromoCodeInput)

	promocodeID := primitive.NewObjectID()
	offerId, _ := primitive.ObjectIDFromHex("6034253f561e5c7cbae6e5f2")

	school := domain.School{
		ID: primitive.NewObjectID(),
		Settings: domain.Settings{
			Domains: []string{"localhost"},
		},
	}

	tests := []struct {
		name         string
		body         string
		school       domain.School
		input        service.CreatePromoCodeInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			body:   `{"code": "TESTPROMO", "discountPercentage": 15, "expiresAt": "2022-12-10T13:49:51.0Z", "offerIds": ["6034253f561e5c7cbae6e5f2"]}`,
			school: school,
			input: service.CreatePromoCodeInput{
				SchoolID:           school.ID,
				Code:               "TESTPROMO",
				DiscountPercentage: 15,
				ExpiresAt:          time.Date(2022, 12, 10, 13, 49, 51, 0, time.UTC),
				OfferIDs:           []primitive.ObjectID{offerId},
			},
			mockBehavior: func(r *mock_service.MockPromoCodes, input service.CreatePromoCodeInput) {
				r.EXPECT().Create(context.Background(), input).Return(promocodeID, nil)
			},
			statusCode:   201,
			responseBody: fmt.Sprintf(`{"id":"%s"}`, promocodeID.Hex()),
		},
		{
			name:         "invalid input body param",
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockPromoCodes, input service.CreatePromoCodeInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "service error",
			body:   `{"code": "TESTPROMO", "discountPercentage": 15, "expiresAt": "2022-12-10T13:49:51.0Z", "offerIds": ["6034253f561e5c7cbae6e5f2"]}`,
			school: school,
			input: service.CreatePromoCodeInput{
				SchoolID:           school.ID,
				Code:               "TESTPROMO",
				DiscountPercentage: 15,
				ExpiresAt:          time.Date(2022, 12, 10, 13, 49, 51, 0, time.UTC),
				OfferIDs:           []primitive.ObjectID{offerId},
			},
			mockBehavior: func(r *mock_service.MockPromoCodes, input service.CreatePromoCodeInput) {
				r.EXPECT().Create(context.Background(), input).Return(primitive.ObjectID{}, errors.New("failed to create promocode"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to create promocode"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			p := mock_service.NewMockPromoCodes(c)
			tt.mockBehavior(p, tt.input)

			services := &service.Services{PromoCodes: p}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.POST("/admins/promocodes/", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminCreatePromocode)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/admins/promocodes/",
				bytes.NewBufferString(tt.body))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

func TestHandler_adminGetPromocodes(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPromoCodes, schoolId primitive.ObjectID)

	type offersMockBehavior func(r *mock_service.MockOffers, offerIds []primitive.ObjectID)

	promocodeId := primitive.NewObjectID()
	offerId := primitive.NewObjectID()

	school := domain.School{
		ID: primitive.NewObjectID(),
		Settings: domain.Settings{
			Domains: []string{"localhost"},
		},
	}

	tests := []struct {
		name               string
		school             domain.School
		mockBehavior       mockBehavior
		offersMockBehavior offersMockBehavior
		statusCode         int
		responseBody       string
	}{
		{
			name:   "ok",
			school: school,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId primitive.ObjectID) {
				r.EXPECT().GetBySchool(context.Background(), schoolId).Return([]domain.PromoCode{
					{
						ID:                 promocodeId,
						SchoolID:           schoolId,
						Code:               "FIRSTPROMO",
						DiscountPercentage: 15,
						ExpiresAt:          time.Date(2022, 12, 10, 13, 49, 51, 0, time.UTC),
						OfferIDs:           []primitive.ObjectID{offerId},
					},
				}, nil)
			},
			offersMockBehavior: func(r *mock_service.MockOffers, offerIds []primitive.ObjectID) {
				r.EXPECT().GetByIds(context.Background(), offerIds).Return([]domain.Offer{
					{
						ID:   offerId,
						Name: "offer",
					},
				}, nil)
			},
			statusCode:   200,
			responseBody: fmt.Sprintf(`{"data":[{"id":"%s","code":"FIRSTPROMO","discountPercentage":15,"expiresAt":"2022-12-10T13:49:51Z","offers":[{"id":"%s","name":"offer"}]}],"count":0}`, promocodeId.Hex(), offerId.Hex()),
		},
		{
			name:   "service error",
			school: school,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId primitive.ObjectID) {
				r.EXPECT().GetBySchool(context.Background(), schoolId).Return(nil, errors.New("failed to get promocodes"))
			},
			offersMockBehavior: func(r *mock_service.MockOffers, offerIds []primitive.ObjectID) {
			},
			statusCode:   500,
			responseBody: `{"message":"failed to get promocodes"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			p := mock_service.NewMockPromoCodes(c)
			tt.mockBehavior(p, tt.school.ID)

			o := mock_service.NewMockOffers(c)
			tt.offersMockBehavior(o, []primitive.ObjectID{offerId})

			services := &service.Services{PromoCodes: p, Offers: o}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/admins/promocodes/", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminGetPromocodes)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/admins/promocodes/", nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

func TestHandler_adminGetPromocodeById(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPromoCodes, schoolId primitive.ObjectID, id primitive.ObjectID)

	type offersMockBehavior func(r *mock_service.MockOffers, offerIds []primitive.ObjectID)

	promocodeId := primitive.NewObjectID()
	offerId := primitive.NewObjectID()

	school := domain.School{
		ID: primitive.NewObjectID(),
		Settings: domain.Settings{
			Domains: []string{"localhost"},
		},
	}

	tests := []struct {
		name               string
		school             domain.School
		mockBehavior       mockBehavior
		offersMockBehavior offersMockBehavior
		statusCode         int
		responseBody       string
	}{
		{
			name:   "ok",
			school: school,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId primitive.ObjectID, id primitive.ObjectID) {
				r.EXPECT().GetById(context.Background(), schoolId, id).Return(domain.PromoCode{
					ID:                 promocodeId,
					SchoolID:           schoolId,
					Code:               "FIRSTPROMO",
					DiscountPercentage: 15,
					ExpiresAt:          time.Date(2022, 12, 10, 13, 49, 51, 0, time.UTC),
					OfferIDs:           []primitive.ObjectID{offerId},
				}, nil)
			},
			offersMockBehavior: func(r *mock_service.MockOffers, offerIds []primitive.ObjectID) {
				r.EXPECT().GetByIds(context.Background(), offerIds).Return([]domain.Offer{
					{
						ID:   offerId,
						Name: "offer",
					},
				}, nil)
			},
			statusCode:   200,
			responseBody: fmt.Sprintf(`{"id":"%s","code":"FIRSTPROMO","discountPercentage":15,"expiresAt":"2022-12-10T13:49:51Z","offers":[{"id":"%s","name":"offer"}]}`, promocodeId.Hex(), offerId.Hex()),
		},
		{
			name:   "service error",
			school: school,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId primitive.ObjectID, id primitive.ObjectID) {
				r.EXPECT().GetById(context.Background(), schoolId, id).Return(domain.PromoCode{}, errors.New("failed to get promocode by id"))
			},
			offersMockBehavior: func(r *mock_service.MockOffers, offerIds []primitive.ObjectID) {
			},
			statusCode:   500,
			responseBody: `{"message":"failed to get promocode by id"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			p := mock_service.NewMockPromoCodes(c)
			tt.mockBehavior(p, tt.school.ID, promocodeId)

			o := mock_service.NewMockOffers(c)
			tt.offersMockBehavior(o, []primitive.ObjectID{offerId})

			services := &service.Services{PromoCodes: p, Offers: o}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.GET("/admins/promocodes/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminGetPromocodeById)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/admins/promocodes/%s", promocodeId.Hex()), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

func TestHandler_adminUpdatePromocode(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPromoCodes, input domain.UpdatePromoCodeInput)

	school := domain.School{
		ID: primitive.NewObjectID(),
		Settings: domain.Settings{
			Domains: []string{"localhost"},
		},
	}

	tests := []struct {
		name         string
		body         string
		school       domain.School
		input        domain.UpdatePromoCodeInput
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			body:   `{"code": "TESTPROMO", "discountPercentage": 15}`,
			school: school,
			input: domain.UpdatePromoCodeInput{
				ID:                 primitive.NewObjectID(),
				SchoolID:           school.ID,
				Code:               "TESTPROMO",
				DiscountPercentage: 15,
			},
			mockBehavior: func(r *mock_service.MockPromoCodes, input domain.UpdatePromoCodeInput) {
				r.EXPECT().Update(context.Background(), input).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:         "invalid input body",
			body:         `{wrong}`,
			school:       school,
			mockBehavior: func(r *mock_service.MockPromoCodes, input domain.UpdatePromoCodeInput) {},
			statusCode:   400,
			responseBody: `{"message":"invalid input body"}`,
		},
		{
			name:   "service error",
			body:   `{"code": "TESTPROMO", "discountPercentage": 15}`,
			school: school,
			input: domain.UpdatePromoCodeInput{
				ID:                 primitive.NewObjectID(),
				SchoolID:           school.ID,
				Code:               "TESTPROMO",
				DiscountPercentage: 15,
			},
			mockBehavior: func(r *mock_service.MockPromoCodes, input domain.UpdatePromoCodeInput) {
				r.EXPECT().Update(context.Background(), input).Return(errors.New("failed to update promocode"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to update promocode"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			p := mock_service.NewMockPromoCodes(c)
			tt.mockBehavior(p, tt.input)

			services := &service.Services{PromoCodes: p}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.PUT("/admins/promocodes/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminUpdatePromocode)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/admins/promocodes/%s", tt.input.ID.Hex()),
				bytes.NewBufferString(tt.body))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

func TestHandler_adminDeletePromocode(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPromoCodes, schoolId, id primitive.ObjectID)

	promocodeId := primitive.NewObjectID()

	school := domain.School{
		ID: primitive.NewObjectID(),
		Settings: domain.Settings{
			Domains: []string{"localhost"},
		},
	}

	tests := []struct {
		name         string
		school       domain.School
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:   "ok",
			school: school,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId, id primitive.ObjectID) {
				r.EXPECT().Delete(context.Background(), schoolId, id).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:   "service error",
			school: school,
			mockBehavior: func(r *mock_service.MockPromoCodes, schoolId, id primitive.ObjectID) {
				r.EXPECT().Delete(context.Background(), schoolId, id).Return(errors.New("failed to delete promocode"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to delete promocode"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			p := mock_service.NewMockPromoCodes(c)
			tt.mockBehavior(p, tt.school.ID, promocodeId)

			services := &service.Services{PromoCodes: p}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.DELETE("/admins/promocodes/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminDeletePromocode)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/admins/promocodes/%s", promocodeId.Hex()), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

func TestHandler_adminDeleteCourse(t *testing.T) {
	type mockBehavior func(r *mock_service.MockCourses, schoolId, id primitive.ObjectID)

	school := domain.School{
		ID: primitive.NewObjectID(),
		Settings: domain.Settings{
			Domains: []string{"localhost"},
		},
	}

	tests := []struct {
		name         string
		courseId     primitive.ObjectID
		school       domain.School
		mockBehavior mockBehavior
		statusCode   int
		responseBody string
	}{
		{
			name:     "ok",
			school:   school,
			courseId: primitive.NewObjectID(),
			mockBehavior: func(r *mock_service.MockCourses, schoolId, id primitive.ObjectID) {
				r.EXPECT().Delete(context.Background(), schoolId, id).Return(nil)
			},
			statusCode:   200,
			responseBody: "",
		},
		{
			name:     "service error",
			school:   school,
			courseId: primitive.NewObjectID(),
			mockBehavior: func(r *mock_service.MockCourses, schoolId, id primitive.ObjectID) {
				r.EXPECT().Delete(context.Background(), schoolId, id).Return(errors.New("failed to delete course"))
			},
			statusCode:   500,
			responseBody: `{"message":"failed to delete course"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			p := mock_service.NewMockCourses(c)
			tt.mockBehavior(p, tt.school.ID, tt.courseId)

			services := &service.Services{Courses: p}
			handler := Handler{services: services}

			// Init Endpoint
			r := gin.New()
			r.DELETE("/admins/courses/:id", func(c *gin.Context) {
				c.Set(schoolCtx, tt.school)
			}, handler.adminDeleteCourse)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/admins/courses/%s", tt.courseId.Hex()), nil)

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tt.statusCode)
			assert.Equal(t, w.Body.String(), tt.responseBody)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
