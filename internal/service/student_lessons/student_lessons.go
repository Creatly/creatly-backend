package studentlessons

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	StudentLessons interface {
		AddFinished(ctx context.Context, studentId, lessonId primitive.ObjectID) error
		SetLastOpened(ctx context.Context, studentId, lessonId primitive.ObjectID) error
	}

	StudentLessonsService struct {
		repo StudentLessons
	}
)

func NewStudentLessonsService(repo StudentLessons) *StudentLessonsService {
	return &StudentLessonsService{
		repo: repo,
	}
}

func (s *StudentLessonsService) AddFinished(ctx context.Context, studentId, lessonId primitive.ObjectID) error {
	return s.repo.AddFinished(ctx, studentId, lessonId)
}

func (s *StudentLessonsService) SetLastOpened(ctx context.Context, studentId, lessonId primitive.ObjectID) error {
	return s.repo.SetLastOpened(ctx, studentId, lessonId)
}
