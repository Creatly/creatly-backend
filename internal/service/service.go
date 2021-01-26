package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/auth"
	"github.com/zhashkevych/courses-backend/pkg/cache"
	"github.com/zhashkevych/courses-backend/pkg/email"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"github.com/zhashkevych/courses-backend/pkg/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

type StudentSignInInput struct {
	Email    string
	Password string
	SchoolID primitive.ObjectID
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Students interface {
	SignUp(ctx context.Context, input StudentSignUpInput) error
	SignIn(ctx context.Context, input StudentSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (Tokens, error)
	Verify(ctx context.Context, hash string) error
	GetStudentModuleWithLessons(ctx context.Context, schoolId, studentId, moduleId primitive.ObjectID) ([]domain.Lesson, error)
	CreateOrder(ctx context.Context, studentId, offerId, promocodeId primitive.ObjectID) (string, error)
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

type Courses interface {
	GetCourseModules(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
	GetModule(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetModuleWithContent(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetModuleOffers(ctx context.Context, schoolId, moduleId primitive.ObjectID) ([]domain.Offer, error)
	GetPackageOffers(ctx context.Context, schoolId, packageId primitive.ObjectID) ([]domain.Offer, error)
	GetPromocodeByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.Promocode, error)
	GetPromocodeById(ctx context.Context, id primitive.ObjectID) (domain.Promocode, error)
	GetOfferById(ctx context.Context, id primitive.ObjectID) (domain.Offer, error)
}

type Services struct {
	Schools  Schools
	Students Students
	Courses  Courses
}

func NewServices(repos *repository.Repositories, cache cache.Cache, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	emailProvider email.Provider, emailListID string, paymentProvider payment.Provider, accessTTL, refreshTTL time.Duration) *Services {

	emailsService := NewEmailsService(emailProvider, emailListID)
	coursesService := NewCoursesService(repos.Courses, repos.Offers, repos.Promocodes)

	return &Services{
		Schools:  NewSchoolsService(repos.Schools, cache),
		Students: NewStudentsService(repos.Students, coursesService, hasher, tokenManager, emailsService, paymentProvider, accessTTL, refreshTTL),
		Courses:  coursesService,
	}
}
