package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
)

type StudentsService struct {
	repo repository.Students
}

func NewStudentsService(repo repository.Students) *StudentsService {
	return &StudentsService{repo: repo}
}

func (s *StudentsService) SignIn(ctx context.Context, email, password string) (string, error) {
	return "", nil
}

func (s *StudentsService) SignUp(ctx context.Context, input StudentSignUpInput) error {
	student := domain.Student{
		Name:     input.Name,
		Password: input.Password, // hash
		Email:    input.Email,
		//SourceCourseID: , str -> oid
		//RegisteredAt: primitive.,
		SchoolID: input.SchoolID,
		Verification: domain.Verification{
			Hash: "", // generate hash
		},
	}

	// TODO send emails

	return s.repo.Create(ctx, student)
}

func (s *StudentsService) Verify(ctx context.Context, hash string) error {
	return nil
}
