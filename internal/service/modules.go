package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModulesService struct {
	repo repository.Courses
}

func NewModulesService(repo repository.Courses) *ModulesService {
	return &ModulesService{repo: repo}
}

func (s *ModulesService) GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	modules, err := s.repo.GetModules(ctx, courseId)
	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (s *ModulesService) GetById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	return s.repo.GetModule(ctx, moduleId)
}

func (s *ModulesService) GetWithContent(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	return s.repo.GetModuleWithContent(ctx, moduleId)
}

func (s *ModulesService) GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error) {
	return s.repo.GetPackagesModules(ctx, packageIds)
}

func (s *ModulesService) Create(ctx context.Context, courseId primitive.ObjectID, name string, position int) (primitive.ObjectID, error) {
	module := domain.Module{
		Name:     name,
		Position: position,
		CourseID: courseId,
	}

	return s.repo.CreateModule(ctx, module)
}
