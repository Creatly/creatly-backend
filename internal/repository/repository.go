package repository

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository/mdb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Schools interface {
	GetByDomain(ctx context.Context, domain string) (domain.School, error)
}

type Students interface {
	Create(ctx context.Context, student domain.Student) error
	GetByCredentials(ctx context.Context, schoolId primitive.ObjectID, email, password string) (domain.Student, error)
	GetByRefreshToken(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (domain.Student, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Student, error)
	SetSession(ctx context.Context, studentId primitive.ObjectID, session domain.Session) error
	GiveModuleAccess(ctx context.Context, studentId, moduleId primitive.ObjectID) error
	Verify(ctx context.Context, code string) error
}

type Courses interface {
	GetModules(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
	GetModule(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetModuleWithContent(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
}

type Offers interface {
	GetSchoolOffers(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Offer, error)
}

type Repositories struct {
	Schools  Schools
	Students Students
	Courses  Courses
	Offers   Offers
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Schools:  mdb.NewSchoolsRepo(db),
		Students: mdb.NewStudentsRepo(db),
		Courses:  mdb.NewCoursesRepo(db),
		Offers:   mdb.NewOffersRepo(db),
	}
}
