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
	GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error)
}

type Students interface {
	Create(ctx context.Context, student domain.Student) error
	GetByCredentials(ctx context.Context, schoolId primitive.ObjectID, email, password string) (domain.Student, error)
	GetByRefreshToken(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (domain.Student, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Student, error)
	SetSession(ctx context.Context, studentId primitive.ObjectID, session domain.Session) error
	GiveAccessToModules(ctx context.Context, studentId primitive.ObjectID, moduleIds []primitive.ObjectID) error
	Verify(ctx context.Context, code string) error
}

type Admins interface {
	GetByCredentials(ctx context.Context, schoolId primitive.ObjectID, email, password string) (domain.Admin, error)
	GetByRefreshToken(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (domain.Admin, error)
	SetSession(ctx context.Context, id primitive.ObjectID, session domain.Session) error
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Admin, error)
}

type Courses interface {
	GetModules(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
	GetModule(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetModuleWithContent(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetPackagesModules(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error)
	Create(ctx context.Context, schoolId primitive.ObjectID, course domain.Course) (primitive.ObjectID, error)
	UpdateCourse(ctx context.Context, schoolId primitive.ObjectID, course domain.Course) error
}

type Offers interface {
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Offer, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Offer, error)
}

type Promocodes interface {
	GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.Promocode, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Promocode, error)
}

type Orders interface {
	Create(ctx context.Context, order domain.Order) error
	AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error)
}

type Repositories struct {
	Schools    Schools
	Students   Students
	Courses    Courses
	Offers     Offers
	Promocodes Promocodes
	Orders     Orders
	Admins     Admins
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Schools:    mdb.NewSchoolsRepo(db),
		Students:   mdb.NewStudentsRepo(db),
		Courses:    mdb.NewCoursesRepo(db),
		Offers:     mdb.NewOffersRepo(db),
		Promocodes: mdb.NewPromocodeRepo(db),
		Orders:     mdb.NewOrdersRepo(db),
		Admins:     mdb.NewAdminsRepo(db),
	}
}
