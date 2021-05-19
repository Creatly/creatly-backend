package service_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_repository "github.com/zhashkevych/creatly-backend/internal/repository/mocks"
	"github.com/zhashkevych/creatly-backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func mockCoursesService(t *testing.T) (*service.CoursesService, *mock_repository.MockCourses, *service.ModulesService, *mock_repository.MockLessonContent) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	coursesRepo := mock_repository.NewMockCourses(mockCtl)
	modulesRepo := mock_repository.NewMockModules(mockCtl)
	lessonsContentRepo := mock_repository.NewMockLessonContent(mockCtl)

	modulesService := service.NewModulesService(modulesRepo, lessonsContentRepo)

	coursesService := service.NewCoursesService(coursesRepo, modulesService)

	return coursesService, coursesRepo, modulesService, lessonsContentRepo
}

func TestNewCoursesService_CreateErr(t *testing.T) {
	coursesService, coursesRepo, _, _ := mockCoursesService(t)

	ctx := context.Background()

	coursesRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(primitive.ObjectID{}, errInternalServErr)

	_, err := coursesService.Create(context.Background(), primitive.ObjectID{}, "")

	require.Error(t, err)
}

func TestNewCoursesService_Create(t *testing.T) {
	coursesService, coursesRepo, _, _ := mockCoursesService(t)

	ctx := context.Background()

	coursesRepo.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(primitive.ObjectID{}, nil)

	_, err := coursesService.Create(context.Background(), primitive.ObjectID{}, "")

	require.NoError(t, err)
}
