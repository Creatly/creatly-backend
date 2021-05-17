package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	mock_repository "github.com/zhashkevych/creatly-backend/internal/repository/mocks"
	"github.com/zhashkevych/creatly-backend/internal/service"
	"github.com/zhashkevych/creatly-backend/pkg/auth"
	"github.com/zhashkevych/creatly-backend/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var errInternalServErr = errors.New("test: internal server error")

func mockAdminService(t *testing.T) (*service.AdminsService, *mock_repository.MockAdmins, *mock_repository.MockSchools) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	adminMock := mock_repository.NewMockAdmins(mockCtl)
	schoolsMock := mock_repository.NewMockSchools(mockCtl)

	adminService := service.NewAdminsService(
		&hash.SHA1Hasher{},
		&auth.Manager{},
		adminMock,
		schoolsMock,
		1*time.Minute,
		1*time.Minute,
	)

	return adminService, adminMock, schoolsMock
}

func TestNewAdminsService_SignInErr(t *testing.T) {
	adminService, adminMock, _ := mockAdminService(t)

	ctx := context.Background()

	adminMock.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.Admin{}, errInternalServErr)
	adminMock.EXPECT().SetSession(ctx, gomock.Any(), gomock.Any())

	res, err := adminService.SignIn(ctx, service.SignInInput{})

	require.True(t, errors.Is(err, errInternalServErr))
	require.Equal(t, service.Tokens{}, res)
}

func TestNewAdminsService_SignIn(t *testing.T) {
	adminService, adminMock, _ := mockAdminService(t)

	ctx := context.Background()

	adminMock.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any(), gomock.Any())
	adminMock.EXPECT().SetSession(ctx, gomock.Any(), gomock.Any())

	res, err := adminService.SignIn(ctx, service.SignInInput{})

	require.NoError(t, err)
	require.IsType(t, service.Tokens{}, res)
}

func TestNewAdminsService_RefreshTokensErr(t *testing.T) {
	adminService, adminMock, _ := mockAdminService(t)

	ctx := context.Background()

	adminMock.EXPECT().GetByRefreshToken(ctx, gomock.Any(), gomock.Any()).Return(domain.Admin{}, errInternalServErr)

	res, err := adminService.RefreshTokens(ctx, primitive.ObjectID{}, "")

	require.True(t, errors.Is(err, errInternalServErr))
	require.Equal(t, service.Tokens{}, res)
}

func TestNewAdminsService_RefreshTokens(t *testing.T) {
	adminService, adminMock, _ := mockAdminService(t)

	ctx := context.Background()

	adminMock.EXPECT().GetByRefreshToken(ctx, gomock.Any(), gomock.Any())
	adminMock.EXPECT().SetSession(ctx, gomock.Any(), gomock.Any())

	res, err := adminService.RefreshTokens(ctx, primitive.ObjectID{}, "")

	require.NoError(t, err)
	require.IsType(t, service.Tokens{}, res)
}

func TestNewAdminsService_GetCoursesErr(t *testing.T) {
	adminService, _, schoolsMock := mockAdminService(t)

	ctx := context.Background()

	schoolsMock.EXPECT().GetById(ctx, gomock.Any()).Return(domain.School{}, errInternalServErr)

	res, err := adminService.GetCourses(ctx, primitive.ObjectID{})

	require.True(t, errors.Is(err, errInternalServErr))
	require.Equal(t, []domain.Course(nil), res)
}

func TestNewAdminsService_GetCourses(t *testing.T) {
	adminService, _, schoolsMock := mockAdminService(t)

	ctx := context.Background()

	schoolsMock.EXPECT().GetById(ctx, gomock.Any())

	res, err := adminService.GetCourses(ctx, primitive.ObjectID{})

	require.NoError(t, err)
	require.IsType(t, []domain.Course{}, res)
}
