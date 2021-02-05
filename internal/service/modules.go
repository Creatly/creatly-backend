package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModulesService struct {
	repo        repository.Modules
	contentRepo repository.LessonContent
}

func NewModulesService(repo repository.Modules, contentRepo repository.LessonContent) *ModulesService {
	return &ModulesService{repo: repo, contentRepo: contentRepo}
}

func (s *ModulesService) GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	modules, err := s.repo.GetByCourse(ctx, courseId)
	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (s *ModulesService) GetById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	return s.repo.GetById(ctx, moduleId)
}

func (s *ModulesService) GetWithContent(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	module, err := s.repo.GetById(ctx, moduleId)
	if err != nil {
		return module, err
	}

	lessonIds := make([]primitive.ObjectID, len(module.Lessons))
	for _, lesson := range module.Lessons {
		lessonIds = append(lessonIds, lesson.ID)
	}

	content, err := s.contentRepo.GetByLessons(ctx, lessonIds)
	if err != nil {
		return module, err
	}

	for i := range module.Lessons {
		for _, lessonContent := range content {
			if module.Lessons[i].ID == lessonContent.LessonID {
				module.Lessons[i].Content = lessonContent.Content
			}
		}
	}

	return module, nil
}

func (s *ModulesService) GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error) {
	return s.repo.GetByPackages(ctx, packageIds)
}

func (s *ModulesService) Create(ctx context.Context, courseId primitive.ObjectID, name string, position int) (primitive.ObjectID, error) {
	module := domain.Module{
		Name:     name,
		Position: position,
		CourseID: courseId,
	}

	return s.repo.Create(ctx, module)
}
