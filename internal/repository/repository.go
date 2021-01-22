package repository

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository/mdb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Schools interface {
	GetByDomain(ctx context.Context, domain string) (domain.School, error)
}

type Students interface {
	Create(ctx context.Context, student domain.Student) error
	GetByCredentials(ctx context.Context, email, password domain.Student) error
	Verify(ctx context.Context, hash string) error
}

type Repositories struct {
	Schools Schools
	Students Students
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Schools: mdb.NewSchoolsRepo(db),
		Students: mdb.NewStudentsRepo(db),
	}
}