package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (s *ModulesService) Create(ctx context.Context, inp CreateModuleInput) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(inp.CourseID)
	if err != nil {
		return id, err
	}

	module := domain.Module{
		Name:     inp.Name,
		Position: inp.Position,
		CourseID: id,
	}

	return s.repo.Create(ctx, module)
}

func (s *ModulesService) Update(ctx context.Context, inp UpdateModuleInput) error {
	id, err := primitive.ObjectIDFromHex(inp.ID)
	if err != nil {
		return err
	}

	updateInput := repository.UpdateModuleInput{
		ID:        id,
		Name:      inp.Name,
		Position:  inp.Position,
		Published: inp.Published,
	}

	return s.repo.Update(ctx, updateInput)
}

func (s *ModulesService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *ModulesService) AddLesson(ctx context.Context, inp AddLessonInput) (primitive.ObjectID, error) {
	lesson := domain.Lesson{
		ID:       primitive.NewObjectID(),
		Name:     inp.Name,
		Position: inp.Position,
	}

	id, err := primitive.ObjectIDFromHex(inp.ModuleID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	err = s.repo.AddLesson(ctx, id, lesson)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return lesson.ID, nil
}

func (s *ModulesService) GetLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.Lesson, error) {
	module, err := s.repo.GetByLesson(ctx, lessonId)
	if err != nil {
		return domain.Lesson{}, err
	}

	var lesson domain.Lesson
	for _, l := range module.Lessons {
		if l.ID == lessonId {
			lesson = l
		}
	}

	content, err := s.contentRepo.GetByLesson(ctx, lessonId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return lesson, nil
		}

		return lesson, err
	}

	lesson.Content = content.Content

	return lesson, nil
}

func (s *ModulesService) UpdateLesson(ctx context.Context, inp UpdateLessonInput) error {
	id, err := primitive.ObjectIDFromHex(inp.LessonID)
	if err != nil {
		return err
	}

	if inp.Name != "" || inp.Position != nil || inp.Published != nil {
		if err := s.repo.UpdateLesson(ctx, repository.UpdateLessonInput{
			ID:        id,
			Name:      inp.Name,
			Position:  inp.Position,
			Published: inp.Published,
		}); err != nil {
			return err
		}
	}

	if inp.Content != "" {
		if err := s.contentRepo.Update(ctx, id, inp.Content); err != nil {
			return err
		}
	}

	return nil
}

func (s *ModulesService) DeleteLesson(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteLesson(ctx, id)
}
