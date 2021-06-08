package lessons

import (
	"context"
	"errors"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Modules interface {
		AddLesson(ctx context.Context, schoolId, id primitive.ObjectID, lesson domain.Lesson) error
		GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.Module, error)
		UpdateLesson(ctx context.Context, inp repository.UpdateLessonInput) error
		DeleteLesson(ctx context.Context, schoolId, id primitive.ObjectID) error
	}

	LessonContent interface {
		GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.LessonContent, error)
		Update(ctx context.Context, schoolId, lessonId primitive.ObjectID, content string) error
		DeleteContent(ctx context.Context, schoolId primitive.ObjectID, lessonIds []primitive.ObjectID) error
	}

	UpdateLessonInput struct {
		LessonID  string
		SchoolID  string
		Name      string
		Content   string
		Position  *uint
		Published *bool
	}

	LessonsService struct {
		repo        Modules
		contentRepo LessonContent
	}

	AddLessonInput struct {
		ModuleID string
		SchoolID string
		Name     string
		Position uint
	}
)

func NewLessonsService(repo Modules, contentRepo LessonContent) *LessonsService {
	return &LessonsService{
		repo:        repo,
		contentRepo: contentRepo,
	}
}

func (s *LessonsService) Create(ctx context.Context, inp AddLessonInput) (primitive.ObjectID, error) {
	schoolId, err := primitive.ObjectIDFromHex(inp.SchoolID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	lesson := domain.Lesson{
		ID:       primitive.NewObjectID(),
		SchoolID: schoolId,
		Name:     inp.Name,
		Position: inp.Position,
	}

	id, err := primitive.ObjectIDFromHex(inp.ModuleID)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	if err := s.repo.AddLesson(ctx, schoolId, id, lesson); err != nil {
		return primitive.ObjectID{}, err
	}

	return lesson.ID, nil
}

func (s *LessonsService) GetById(ctx context.Context, lessonId primitive.ObjectID) (domain.Lesson, error) {
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return lesson, nil
		}

		return lesson, err
	}

	lesson.Content = content.Content

	return lesson, nil
}

func (s *LessonsService) Update(ctx context.Context, inp UpdateLessonInput) error {
	id, err := primitive.ObjectIDFromHex(inp.LessonID)
	if err != nil {
		return err
	}

	schoolId, err := primitive.ObjectIDFromHex(inp.SchoolID)
	if err != nil {
		return err
	}

	if inp.Name != "" || inp.Position != nil || inp.Published != nil {
		if err := s.repo.UpdateLesson(ctx, repository.UpdateLessonInput{
			ID:        id,
			Name:      inp.Name,
			Position:  inp.Position,
			Published: inp.Published,
			SchoolID:  schoolId,
		}); err != nil {
			return err
		}
	}

	if inp.Content != "" {
		if err := s.contentRepo.Update(ctx, schoolId, id, inp.Content); err != nil {
			return err
		}
	}

	return nil
}

func (s *LessonsService) Delete(ctx context.Context, schoolId, id primitive.ObjectID) error {
	return s.repo.DeleteLesson(ctx, schoolId, id)
}

func (s *LessonsService) DeleteContent(ctx context.Context, schoolId primitive.ObjectID, lessonIds []primitive.ObjectID) error {
	return s.contentRepo.DeleteContent(ctx, schoolId, lessonIds)
}
