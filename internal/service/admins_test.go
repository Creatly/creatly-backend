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
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	adminRepo := mock_repository.NewMockAdmins(mockCtl)
	schoolsRepo := mock_repository.NewMockSchools(mockCtl)
	studentsRepo := mock_repository.NewMockStudents(mockCtl)

	adminService := service.NewAdminsService(
		&hash.SHA1Hasher{},
		&auth.Manager{},
		adminRepo,
		schoolsRepo,
		studentsRepo,
		1*time.Minute,
		1*time.Minute,
	)

	return adminService, adminRepo, schoolsRepo
}

func TestNewAdminsService_SignInErr(t *testing.T) {
	adminService, adminRepo, _ := mockAdminService(t)

	ctx := context.Background()

	adminRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(domain.Admin{}, errInternalServErr)
	adminRepo.EXPECT().SetSession(ctx, gomock.Any(), gomock.Any())

	res, err := adminService.SignIn(ctx, service.SchoolSignInInput{})

	require.True(t, errors.Is(err, errInternalServErr))
	require.Equal(t, service.Tokens{}, res)
}

func TestNewAdminsService_SignIn(t *testing.T) {
	adminService, adminRepo, _ := mockAdminService(t)

	ctx := context.Background()

	adminRepo.EXPECT().GetByCredentials(ctx, gomock.Any(), gomock.Any(), gomock.Any())
	adminRepo.EXPECT().SetSession(ctx, gomock.Any(), gomock.Any())

	res, err := adminService.SignIn(ctx, service.SchoolSignInInput{})

	require.NoError(t, err)
	require.IsType(t, service.Tokens{}, res)
}

func TestNewAdminsService_RefreshTokensErr(t *testing.T) {
	adminService, adminRepo, _ := mockAdminService(t)

	ctx := context.Background()

	adminRepo.EXPECT().GetByRefreshToken(ctx, gomock.Any(), gomock.Any()).Return(domain.Admin{}, errInternalServErr)

	res, err := adminService.RefreshTokens(ctx, primitive.ObjectID{}, "")

	require.True(t, errors.Is(err, errInternalServErr))
	require.Equal(t, service.Tokens{}, res)
}

func TestNewAdminsService_RefreshTokens(t *testing.T) {
	adminService, adminRepo, _ := mockAdminService(t)

	ctx := context.Background()

	adminRepo.EXPECT().GetByRefreshToken(ctx, gomock.Any(), gomock.Any())
	adminRepo.EXPECT().SetSession(ctx, gomock.Any(), gomock.Any())

	res, err := adminService.RefreshTokens(ctx, primitive.ObjectID{}, "")

	require.NoError(t, err)
	require.IsType(t, service.Tokens{}, res)
}

func TestNewAdminsService_GetCoursesErr(t *testing.T) {
	adminService, _, schoolsRepo := mockAdminService(t)

	ctx := context.Background()

	schoolsRepo.EXPECT().GetById(ctx, gomock.Any()).Return(domain.School{}, errInternalServErr)

	res, err := adminService.GetCourses(ctx, primitive.ObjectID{})

	require.True(t, errors.Is(err, errInternalServErr))
	require.Equal(t, []domain.Course(nil), res)
}

func TestNewAdminsService_GetCourses(t *testing.T) {
	adminService, _, schoolsRepo := mockAdminService(t)

	ctx := context.Background()

	schoolsRepo.EXPECT().GetById(ctx, gomock.Any())

	res, err := adminService.GetCourses(ctx, primitive.ObjectID{})

	require.NoError(t, err)
	require.IsType(t, []domain.Course{}, res)
}

func TestNewAdminsService_GetCourseByIdErr(t *testing.T) {
	adminService, _, schoolsRepo := mockAdminService(t)

	ctx := context.Background()

	schoolsRepo.EXPECT().GetById(ctx, gomock.Any()).Return(domain.School{}, errInternalServErr)

	res, err := adminService.GetCourseById(ctx, primitive.ObjectID{}, primitive.ObjectID{})

	require.True(t, errors.Is(err, errInternalServErr))
	require.Equal(t, domain.Course{}, res)
}

func TestNewAdminsService_GetCourseByIdNotFoundErr(t *testing.T) {
	adminService, _, schoolsRepo := mockAdminService(t)

	ctx := context.Background()

	schoolsRepo.EXPECT().GetById(ctx, gomock.Any())

	_, err := adminService.GetCourseById(ctx, primitive.ObjectID{}, primitive.ObjectID{})

	require.Error(t, err)
}

func TestNewAdminsService_GetCourseById(t *testing.T) {
	adminService, _, schoolsRepo := mockAdminService(t)

	ctx := context.Background()
	s := domain.School{
		ID: primitive.NewObjectID(),
		Courses: []domain.Course{
			{
				ID: primitive.NewObjectID(),
			},
		},
	}

	schoolsRepo.EXPECT().GetById(ctx, gomock.Any()).Return(s, nil)

	_, err := adminService.GetCourseById(ctx, s.ID, s.Courses[0].ID)

	require.NoError(t, err)
}
