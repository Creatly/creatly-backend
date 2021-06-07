package service

import (
	"context"
	"errors"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
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
