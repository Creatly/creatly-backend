package modules

import (
	"context"
	"sort"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CreateModuleInput struct {
		SchoolID string
		CourseID string
		Name     string
		Position uint
	}

	UpdateModuleInput struct {
		ID        string
		SchoolID  string
		Name      string
		Position  *uint
		Published *bool
	}

	UpdateLessonInput struct {
		ID        primitive.ObjectID
		SchoolID  primitive.ObjectID
		Name      string
		Position  *uint
		Published *bool
	}

	Modules interface {
		Create(ctx context.Context, module domain.Module) (primitive.ObjectID, error)
		GetPublishedByCourseId(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
		GetByCourseId(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
		GetPublishedById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
		GetById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
		GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error)
		Update(ctx context.Context, inp repository.UpdateModuleInput) error
		Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
		DeleteByCourse(ctx context.Context, schoolId, courseId primitive.ObjectID) error
		AddLesson(ctx context.Context, schoolId, id primitive.ObjectID, lesson domain.Lesson) error
		GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.Module, error)
		UpdateLesson(ctx context.Context, inp UpdateLessonInput) error
		DeleteLesson(ctx context.Context, schoolId, id primitive.ObjectID) error
		AttachPackage(ctx context.Context, schoolId, packageId primitive.ObjectID, modules []primitive.ObjectID) error
	}

	LessonContent interface {
		GetByLessons(ctx context.Context, lessonIds []primitive.ObjectID) ([]domain.LessonContent, error)
		GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.LessonContent, error)
		Update(ctx context.Context, schoolId, lessonId primitive.ObjectID, content string) error
		DeleteContent(ctx context.Context, schoolId primitive.ObjectID, lessonIds []primitive.ObjectID) error
	}

	ModulesService struct {
		repo        Modules
		contentRepo LessonContent
	}
)

func NewModulesService(repo Modules, contentRepo repository.LessonContent) *ModulesService {
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

func (s *ModulesService) GetWithContent(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	module, err := s.repo.GetById(ctx, moduleId)
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

func (s *ModulesService) GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.Module, error) {
	return s.repo.GetByLesson(ctx, lessonId)
}

func (s *ModulesService) Create(ctx context.Context, inp CreateModuleInput) (primitive.ObjectID, error) {
	id, err := primitive.ObjectIDFromHex(inp.CourseID)
	if err != nil {
		return id, err
	}

	schoolId, err := primitive.ObjectIDFromHex(inp.SchoolID)
	if err != nil {
		return id, err
	}

	module := domain.Module{
		Name:     inp.Name,
		Position: inp.Position,
		CourseID: id,
		SchoolID: schoolId,
	}

	return s.repo.Create(ctx, module)
}

func (s *ModulesService) Update(ctx context.Context, inp UpdateModuleInput) error {
	id, err := primitive.ObjectIDFromHex(inp.ID)
	if err != nil {
		return err
	}

	schoolId, err := primitive.ObjectIDFromHex(inp.SchoolID)
	if err != nil {
		return err
	}

	updateInput := repository.UpdateModuleInput{
		ID:        id,
		SchoolID:  schoolId,
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
