package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_repository "github.com/zhashkevych/creatly-backend/internal/repository/mocks"
	"github.com/zhashkevych/creatly-backend/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	testCoursesServiceUpdate struct {
		input       service.UpdateCourseInput
		mock        func()
		expectedErr error
	}
	testCoursesServiceDelete struct {
		schoolID    primitive.ObjectID
		courseID    primitive.ObjectID
		mock        func()
		expectedErr error
	}
)

func mockCoursesService(t *testing.T) (*service.CoursesService, *mock_repository.MockCourses, *mock_repository.MockModules, *mock_repository.MockLessonContent) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	coursesRepo := mock_repository.NewMockCourses(mockCtl)
	modulesRepo := mock_repository.NewMockModules(mockCtl)
	lessonsContentRepo := mock_repository.NewMockLessonContent(mockCtl)

	modulesService := service.NewModulesService(modulesRepo, lessonsContentRepo)

	coursesService := service.NewCoursesService(coursesRepo, modulesService)

	return coursesService, coursesRepo, modulesRepo, lessonsContentRepo
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

func TestCoursesServiceUpdate(t *testing.T) {
	t.Parallel()

	coursesService, coursesRepo, _, _ := mockCoursesService(t)
	ctx := context.Background()

	tests := map[string]testCoursesServiceUpdate{
		"invalid courses id": {
			input: service.UpdateCourseInput{
				CourseID: "",
			},
			mock:        func() {},
			expectedErr: primitive.ErrInvalidHex,
		},
		"invalid school id": {
			input: service.UpdateCourseInput{
				CourseID: primitive.NewObjectID().Hex(),
				SchoolID: "",
			},
			mock:        func() {},
			expectedErr: primitive.ErrInvalidHex,
		},
		"update course": {
			input: service.UpdateCourseInput{
				CourseID: primitive.NewObjectID().Hex(),
				SchoolID: primitive.NewObjectID().Hex(),
			},
			mock: func() {
				coursesRepo.EXPECT().Update(ctx, gomock.Any())
			},
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			err := coursesService.Update(ctx, tc.input)

			require.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}

func TestCoursesServiceDelete(t *testing.T) {
	t.Parallel()

	coursesService, coursesRepo, modulesRepo, lessonsContentRepo := mockCoursesService(t)
	ctx := context.Background()

	tests := map[string]testCoursesServiceDelete{
		"delete course err": {
			schoolID: primitive.NewObjectID(),
			courseID: primitive.NewObjectID(),
			mock: func() {
				coursesRepo.EXPECT().Delete(ctx, gomock.Any(), gomock.Any()).Return(errInternalServErr)
			},
			expectedErr: errInternalServErr,
		},
		"delete module by course err": {
			schoolID: primitive.NewObjectID(),
			courseID: primitive.NewObjectID(),
			mock: func() {
				coursesRepo.EXPECT().Delete(ctx, gomock.Any(), gomock.Any())
				modulesRepo.EXPECT().GetPublishedByCourseId(ctx, gomock.Any()).Return(nil, errInternalServErr)
				modulesRepo.EXPECT().DeleteByCourse(ctx, gomock.Any(), gomock.Any())
				lessonsContentRepo.EXPECT().DeleteContent(ctx, gomock.Any(), gomock.Any())
				modulesRepo.EXPECT().Delete(ctx, gomock.Any(), gomock.Any())
			},
			expectedErr: errInternalServErr,
		},
		"delete module by course": {
			schoolID: primitive.NewObjectID(),
			courseID: primitive.NewObjectID(),
			mock: func() {
				coursesRepo.EXPECT().Delete(ctx, gomock.Any(), gomock.Any())
				modulesRepo.EXPECT().GetPublishedByCourseId(ctx, gomock.Any())
				modulesRepo.EXPECT().DeleteByCourse(ctx, gomock.Any(), gomock.Any())
				lessonsContentRepo.EXPECT().DeleteContent(ctx, gomock.Any(), gomock.Any())
				modulesRepo.EXPECT().Delete(ctx, gomock.Any(), gomock.Any())
			},
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			err := coursesService.Delete(ctx, tc.schoolID, tc.courseID)

			require.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}
