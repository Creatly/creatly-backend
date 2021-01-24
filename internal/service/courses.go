package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CoursesService struct {
	repo repository.Courses
}

func NewCoursesService(repo repository.Courses) *CoursesService {
	return &CoursesService{repo: repo}
}

func (s *CoursesService) GetCourseModules(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	modules, err := s.repo.GetModules(ctx, courseId)
	if err != nil {
		return nil, err
	}

	if len(modules) == 0 {
		return nil, ErrCourseContentNotFound
	}

	return modules, nil
}
