package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/cache"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schools interface {
	GetByDomain(ctx context.Context, domainName string) (domain.School, error)
}

type StudentSignUpInput struct {
	Name           string
	Email          string
	Password       string
	SchoolID       primitive.ObjectID
	SourceCourseID string
}

type Students interface {
	SignIn(ctx context.Context, email, password string) (string, error)
	SignUp(ctx context.Context, input StudentSignUpInput) error
	Verify(ctx context.Context, hash string) error
}

type Services struct {
	Schools  Schools
	Students Students
}

func NewServices(repos *repository.Repositories, cache cache.Cache) *Services {
	return &Services{
		Schools:  NewSchoolsService(repos.Schools, cache),
		Students: NewStudentsService(repos.Students),
	}
}
