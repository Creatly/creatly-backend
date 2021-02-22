package service

import (
	"context"

	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentLessonsService struct {
	repo repository.StudentLessons
}

func NewStudentLessonsService(repo repository.StudentLessons) *StudentLessonsService {
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
