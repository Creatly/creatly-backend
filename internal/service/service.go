package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/cache"
	"github.com/zhashkevych/courses-backend/pkg/email"
	"github.com/zhashkevych/courses-backend/pkg/hash"
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
	RegisterSource string
}

type SignInResult struct {
	AccessToken string
	RefreshToken string
}

type Students interface {
	SignIn(ctx context.Context, email, password string) (string, error)
	SignUp(ctx context.Context, input StudentSignUpInput) error
	Verify(ctx context.Context, hash string) error
}

type AddToListInput struct {
	Email            string
	Name             string
	RegisterSource   string
	VerificationCode string
}

type Emails interface {
	AddToList(AddToListInput) error
}

type Services struct {
	Schools  Schools
	Students Students
}

func NewServices(repos *repository.Repositories, cache cache.Cache, hasher hash.PasswordHasher, emailProvider email.Provider, emailListID string) *Services {
	emailsService := NewEmailsService(emailProvider, emailListID)
	return &Services{
		Schools:  NewSchoolsService(repos.Schools, cache),
		Students: NewStudentsService(repos.Students, hasher, emailsService),
	}
}
