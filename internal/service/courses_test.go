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

	coursesMock := mock_repository.NewMockCourses(mockCtl)
	modulesMock := mock_repository.NewMockModules(mockCtl)
	lessonsContentMock := mock_repository.NewMockLessonContent(mockCtl)

	modulesService := service.NewModulesService(modulesMock, lessonsContentMock)

	coursesService := service.NewCoursesService(coursesMock, modulesService)

	return coursesService, coursesMock, modulesService, lessonsContentMock
}

func TestNewCoursesService_CreateErr(t *testing.T) {
	coursesService, coursesMock, _, _ := mockCoursesService(t)

	ctx := context.Background()

	coursesMock.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(primitive.ObjectID{}, errInternalServErr)

	_, err := coursesService.Create(context.Background(), primitive.ObjectID{}, "")

	require.Error(t, err)
}

func TestNewCoursesService_Create(t *testing.T) {
	coursesService, coursesMock, _, _ := mockCoursesService(t)

	ctx := context.Background()

	coursesMock.EXPECT().Create(ctx, gomock.Any(), gomock.Any()).Return(primitive.ObjectID{}, nil)

	_, err := coursesService.Create(context.Background(), primitive.ObjectID{}, "")

	require.NoError(t, err)
}
