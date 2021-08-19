package service

import (
	"context"
	"sort"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModulesService struct {
	repo        repository.Modules
	contentRepo repository.LessonContent
}

func NewModulesService(repo repository.Modules, contentRepo repository.LessonContent) *ModulesService {
	return &ModulesService{repo: repo, contentRepo: contentRepo}
}

func (s *ModulesService) GetPublishedByCourseId(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	modules, err := s.repo.GetPublishedByCourseId(ctx, courseId)
	if err != nil {
		return nil, err
	}

	for i := range modules {
		sortLessons(modules[i].Lessons)
	}

	return modules, nil
}

func (s *ModulesService) GetByCourseId(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	modules, err := s.repo.GetByCourseId(ctx, courseId)
	if err != nil {
		return nil, err
	}

	for i := range modules {
		sortLessons(modules[i].Lessons)
	}

	return modules, nil
}

func (s *ModulesService) GetById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	module, err := s.repo.GetPublishedById(ctx, moduleId)
	if err != nil {
		return module, err
	}

	sortLessons(module.Lessons)

	return module, nil
}

func (s *ModulesService) GetWithContent(ctx context.Context, moduleID primitive.ObjectID) (domain.Module, error) {
	module, err := s.repo.GetById(ctx, moduleID)
	if err != nil {
		return module, err
	}

	lessonIds := make([]primitive.ObjectID, len(module.Lessons))
	publishedLessons := make([]domain.Lesson, 0)

	for _, lesson := range module.Lessons {
		if lesson.Published {
			publishedLessons = append(publishedLessons, lesson)
			lessonIds = append(lessonIds, lesson.ID)
		}
	}

	module.Lessons = publishedLessons // remove unpublished lessons from final result

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

	sortLessons(module.Lessons)

	return module, nil
}

func (s *ModulesService) GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error) {
	modules, err := s.repo.GetByPackages(ctx, packageIds)
	if err != nil {
		return nil, err
	}

	for i := range modules {
		sortLessons(modules[i].Lessons)
	}

	return modules, nil
}

func (s *ModulesService) GetByLesson(ctx context.Context, lessonID primitive.ObjectID) (domain.Module, error) {
	return s.repo.GetByLesson(ctx, lessonID)
}

func (s *ModulesService) Create(ctx context.Context, inp CreateModuleInput) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(inp.CourseID)
	if err != nil {
		return id, err
	}

	schoolID, err := primitive.ObjectIDFromHex(inp.SchoolID)
	if err != nil {
		return id, err
	}

	module := domain.Module{
		Name:     inp.Name,
		Position: inp.Position,
		CourseID: id,
		SchoolID: schoolID,
	}

	return s.repo.Create(ctx, module)
}

func (s *ModulesService) Update(ctx context.Context, inp UpdateModuleInput) error {
	id, err := primitive.ObjectIDFromHex(inp.ID)
	if err != nil {
		return err
	}

	schoolID, err := primitive.ObjectIDFromHex(inp.SchoolID)
	if err != nil {
		return err
	}

	updateInput := repository.UpdateModuleInput{
		ID:        id,
		SchoolID:  schoolID,
		Name:      inp.Name,
		Position:  inp.Position,
		Published: inp.Published,
	}

	return s.repo.Update(ctx, updateInput)
}

func (s *ModulesService) Delete(ctx context.Context, schoolId, moduleId primitive.ObjectID) error {
	module, err := s.repo.GetById(ctx, moduleId)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, schoolId, moduleId); err != nil {
		return err
	}

	lessonIds := make([]primitive.ObjectID, len(module.Lessons))
	for _, lesson := range module.Lessons {
		lessonIds = append(lessonIds, lesson.ID)
	}

	return s.contentRepo.DeleteContent(ctx, schoolId, lessonIds)
}

func (s *ModulesService) DeleteByCourse(ctx context.Context, schoolId, courseId primitive.ObjectID) error {
	modules, err := s.repo.GetPublishedByCourseId(ctx, courseId)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteByCourse(ctx, schoolId, courseId); err != nil {
		return err
	}

	lessonIds := make([]primitive.ObjectID, 0)

	for _, module := range modules {
		for _, lesson := range module.Lessons {
			lessonIds = append(lessonIds, lesson.ID)
		}
	}

	return s.contentRepo.DeleteContent(ctx, schoolId, lessonIds)
}

func sortLessons(lessons []domain.Lesson) {
	sort.Slice(lessons, func(i, j int) bool {
		return lessons[i].Position < lessons[j].Position
	})
}
