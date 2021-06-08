package lessons_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mock_repository "github.com/zhashkevych/creatly-backend/internal/repository/mocks"
	"github.com/zhashkevych/creatly-backend/internal/service/lessons"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var errInternalServErr = errors.New("test: internal server error")

type testLessonsServiceCreate struct {
	input       lessons.AddLessonInput
	mock        func()
	expectedErr error
}

func mockLessonsService(t *testing.T) (*lessons.LessonsService, *mock_repository.MockModules, *mock_repository.MockLessonContent) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	modulesRepo := mock_repository.NewMockModules(mockCtl)
	lessonsContentRepo := mock_repository.NewMockLessonContent(mockCtl)

	lessonsService := lessons.NewLessonsService(modulesRepo, lessonsContentRepo)

	return lessonsService, modulesRepo, lessonsContentRepo
}

func TestLessonsServiceCreate(t *testing.T) {
	t.Parallel()

	lessonsService, modulesRepo, _ := mockLessonsService(t)

	ctx := context.Background()

	schoolID := primitive.NewObjectID().Hex()
	moduleID := primitive.NewObjectID().Hex()

	tests := map[string]testLessonsServiceCreate{
		"invalid school id": {
			input: lessons.AddLessonInput{
				SchoolID: "",
				ModuleID: moduleID,
			},
			mock:        func() {},
			expectedErr: primitive.ErrInvalidHex,
		},
		"invalid module id": {
			input: lessons.AddLessonInput{
				SchoolID: schoolID,
				ModuleID: "",
			},
			mock:        func() {},
			expectedErr: primitive.ErrInvalidHex,
		},
		"add lesson error": {
			input: lessons.AddLessonInput{
				SchoolID: schoolID,
				ModuleID: moduleID,
			},
			mock: func() {
				modulesRepo.EXPECT().AddLesson(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(errInternalServErr)
			},
			expectedErr: errInternalServErr,
		},
		"create lesson": {
			input: lessons.AddLessonInput{
				SchoolID: schoolID,
				ModuleID: moduleID,
			},
			mock: func() {
				modulesRepo.EXPECT().AddLesson(ctx, gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tc.mock()

			_, err := lessonsService.Create(ctx, tc.input)

			require.True(t, errors.Is(err, tc.expectedErr))
		})
	}
}
