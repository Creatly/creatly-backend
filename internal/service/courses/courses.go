package courses

import (
	"context"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UpdateCourseInput struct {
		CourseID    string
		SchoolID    string
		Name        string
		Code        string
		Description string
		Color       string
		Published   *bool
	}

	Courses interface {
		Create(ctx context.Context, schoolId primitive.ObjectID, course domain.Course) (primitive.ObjectID, error)
		Update(ctx context.Context, inp repository.UpdateCourseInput) error
		Delete(ctx context.Context, schoolId, courseId primitive.ObjectID) error
	}

	Modules interface {
		DeleteByCourse(ctx context.Context, schoolId, courseId primitive.ObjectID) error
	}

	CoursesService struct {
		repo           Courses
		modulesService Modules
	}
)

func NewCoursesService(repo Courses, modulesService Modules) *CoursesService {
	return &CoursesService{
		repo:           repo,
		modulesService: modulesService,
	}
}

func (s *CoursesService) Create(ctx context.Context, schoolId primitive.ObjectID, name string) (primitive.ObjectID, error) {
	return s.repo.Create(ctx, schoolId, domain.Course{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func (s *CoursesService) Update(ctx context.Context, inp UpdateCourseInput) error {
	updateInput := repository.UpdateCourseInput{
		Name:        inp.Name,
		Code:        inp.Code,
		Description: inp.Description,
		Color:       inp.Color,
		Published:   inp.Published,
	}

	var err error

	updateInput.ID, err = primitive.ObjectIDFromHex(inp.CourseID)
	if err != nil {
		return err
	}

	updateInput.SchoolID, err = primitive.ObjectIDFromHex(inp.SchoolID)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, updateInput)
}

func (s *CoursesService) Delete(ctx context.Context, schoolId, courseId primitive.ObjectID) error {
	if err := s.repo.Delete(ctx, schoolId, courseId); err != nil {
		return err
	}

	return s.modulesService.DeleteByCourse(ctx, schoolId, courseId)
}
