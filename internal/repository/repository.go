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
	SetSession(ctx context.Context, userId primitive.ObjectID, session domain.Session) error
	Verify(ctx context.Context, code string) error
}

type Courses interface {
	GetModules(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
}

type Repositories struct {
	Schools  Schools
	Students Students
	Courses  Courses
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Schools:  mdb.NewSchoolsRepo(db),
		Students: mdb.NewStudentsRepo(db),
		Courses:  mdb.NewCoursesRepo(db),
	}
}
