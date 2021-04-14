package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LessonsService struct {
	repo        repository.Modules
	contentRepo repository.LessonContent
}

func NewLessonsService(repo repository.Modules, contentRepo repository.LessonContent) *LessonsService {
	return &LessonsService{repo: repo, contentRepo: contentRepo}
}

func (s *LessonsService) Create(ctx context.Context, inp AddLessonInput) (primitive.ObjectID, error) {
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
		if err == mongo.ErrNoDocuments {
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

func (s *LessonsService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.DeleteLesson(ctx, id)
}

func (s *LessonsService) DeleteContent(ctx context.Context, lessonIds []primitive.ObjectID) error {
	return s.contentRepo.DeleteContent(ctx, lessonIds)
}
