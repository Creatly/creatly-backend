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

type testStudentLessonsService struct {
	studentID   primitive.ObjectID
	lessonID    primitive.ObjectID
	mock        func()
	expectedErr error
}

func mockStudentLessonsService(t *testing.T) (*service.StudentLessonsService, *mock_repository.MockStudentLessons) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	studentLessonsRepo := mock_repository.NewMockStudentLessons(mockCtl)

	studentLessonsService := service.NewStudentLessonsService(studentLessonsRepo)

	return studentLessonsService, studentLessonsRepo
}
//
//func TestStudentLessonsServiceAddFinished(t *testing.T) {
//	t.Parallel()
//
//	studentLessonsService, studentLessonsRepo := mockStudentLessonsService(t)
//
//	ctx := context.Background()
//
//	studentID := primitive.NewObjectID()
//	lessonID := primitive.NewObjectID()
//
//	tests := map[string]testStudentLessonsService{
//		"repository error": {
//			studentID: studentID,
//			lessonID:  lessonID,
//			mock: func() {
//				studentLessonsRepo.EXPECT().AddFinished(ctx, studentID, lessonID).Return(errInternalServErr)
//			},
//			expectedErr: errInternalServErr,
//		},
//		"add finished": {
//			studentID: studentID,
//			lessonID:  lessonID,
//			mock: func() {
//				studentLessonsRepo.EXPECT().AddFinished(ctx, studentID, lessonID)
//			},
//			expectedErr: nil,
//		},
//	}
//
//	for name, tc := range tests {
//		tc := tc
//
//		t.Run(name, func(t *testing.T) {
//			t.Parallel()
//
//			tc.mock()
//
//			err := studentLessonsService.AddFinished(ctx, tc.studentID, tc.lessonID)
//
//			require.True(t, errors.Is(err, tc.expectedErr))
//		})
//	}
//}

func TestStudentLessonsServiceSetLastOpened(t *testing.T) {
	t.Parallel()

	studentLessonsService, studentLessonsRepo := mockStudentLessonsService(t)

	ctx := context.Background()

	studentID := primitive.NewObjectID()
	lessonID := primitive.NewObjectID()

	tests := map[string]testStudentLessonsService{
		"repository error": {
			studentID: studentID,
			lessonID:  lessonID,
			mock: func() {
				studentLessonsRepo.EXPECT().SetLastOpened(ctx, studentID, lessonID).Return(errInternalServErr)
			},
			expectedErr: errInternalServErr,
		},
		"add finished": {
			studentID: studentID,
			lessonID:  lessonID,
			mock: func() {
				studentLessonsRepo.EXPECT().SetLastOpened(ctx, studentID, lessonID)
			},
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			err := studentLessonsService.SetLastOpened(ctx, tc.studentID, tc.lessonID)

			require.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}
