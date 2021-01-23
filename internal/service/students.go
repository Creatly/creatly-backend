package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type StudentsService struct {
	repo         repository.Students
	hasher       hash.PasswordHasher
	emailService Emails
}

func NewStudentsService(repo repository.Students, hasher hash.PasswordHasher, emailService Emails) *StudentsService {
	return &StudentsService{repo: repo, hasher: hasher, emailService: emailService}
}

func (s *StudentsService) SignIn(ctx context.Context, email, password string) (string, error) {
	return "", nil
}

func (s *StudentsService) SignUp(ctx context.Context, input StudentSignUpInput) error {
	verificationHash := primitive.NewObjectID()
	student := domain.Student{
		Name:           input.Name,
		Password:       s.hasher.Hash(input.Password),
		Email:          input.Email,
		RegisteredAt:   time.Now(),
		LastVisitAt:    time.Now(),
		SchoolID:       input.SchoolID,
		RegisterSource: input.RegisterSource,
		Verification: domain.Verification{
			Hash: verificationHash,
		},
	}

	if err := s.repo.Create(ctx, student); err != nil {
		return err
	}

	// TODO: If it fails, what then?
	return s.emailService.AddToList(AddToListInput{
		Email:            student.Email,
		Name:             student.Name,
		RegisterSource:   student.RegisterSource,
		VerificationCode: verificationHash.Hex(),
	})
}

func (s *StudentsService) Verify(ctx context.Context, hash string) error {
	return s.repo.Verify(ctx, hash)
}
