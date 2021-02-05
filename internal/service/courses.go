package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CoursesService struct {
	repo repository.Courses
}

func NewCoursesService(repo repository.Courses) *CoursesService {
	return &CoursesService{repo: repo}
}

func (s *CoursesService) Create(ctx context.Context, schoolId primitive.ObjectID, name string) (primitive.ObjectID, error) {
	return s.repo.Create(ctx, schoolId, domain.Course{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func (s *CoursesService) Update(ctx context.Context, schoolId primitive.ObjectID, inp UpdateCourseInput) error {
	updateInput := repository.UpdateCourseInput{
		Name:        inp.Name,
		Code:        inp.Code,
		Description: inp.Description,
		Published:   inp.Published,
	}

	var err error
	updateInput.ID, err = primitive.ObjectIDFromHex(inp.CourseID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, schoolId, updateInput)
}
