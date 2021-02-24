package repository

import (
	"context"

	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateSchoolSettingsInput struct {
	SchoolID    primitive.ObjectID
	Color       string
	Domain      string
	Email       string
	ContactData string
	Pages       *domain.Pages
}

type Schools interface {
	GetByDomain(ctx context.Context, domain string) (domain.School, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error)
	UpdateSettings(ctx context.Context, inp UpdateSchoolSettingsInput) error
}

type Students interface {
	Create(ctx context.Context, student domain.Student) error
	GetByCredentials(ctx context.Context, schoolId primitive.ObjectID, email, password string) (domain.Student, error)
	GetByRefreshToken(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (domain.Student, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Student, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Student, error)
	SetSession(ctx context.Context, studentId primitive.ObjectID, session domain.Session) error
	GiveAccessToCourseAndModule(ctx context.Context, studentId, courseId, moduleId primitive.ObjectID) error
	GiveAccessToCoursesAndModules(ctx context.Context, studentId primitive.ObjectID, courseIds, moduleIds []primitive.ObjectID) error
	Verify(ctx context.Context, code string) error
}

type StudentLessons interface {
	AddFinished(ctx context.Context, studentId, lessonId primitive.ObjectID) error
	SetLastOpened(ctx context.Context, studentId, lessonId primitive.ObjectID) error
}

type Admins interface {
	GetByCredentials(ctx context.Context, schoolId primitive.ObjectID, email, password string) (domain.Admin, error)
	GetByRefreshToken(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (domain.Admin, error)
	SetSession(ctx context.Context, id primitive.ObjectID, session domain.Session) error
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Admin, error)
}

type UpdateCourseInput struct {
	ID          primitive.ObjectID
	Name        string
	Code        string
	Description string
	Color       string
	Published   *bool
}

type Courses interface {
	Create(ctx context.Context, schoolId primitive.ObjectID, course domain.Course) (primitive.ObjectID, error)
	Update(ctx context.Context, schoolId primitive.ObjectID, inp UpdateCourseInput) error
}

type UpdateModuleInput struct {
	ID        primitive.ObjectID
	Name      string
	Position  *uint
	Published *bool
}

type UpdateLessonInput struct {
	ID        primitive.ObjectID
	Name      string
	Position  *uint
	Published *bool
}

type Modules interface {
	Create(ctx context.Context, module domain.Module) (primitive.ObjectID, error)
	GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
	GetById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error)
	Update(ctx context.Context, inp UpdateModuleInput) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	AddLesson(ctx context.Context, id primitive.ObjectID, lesson domain.Lesson) error
	GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.Module, error)
	UpdateLesson(ctx context.Context, inp UpdateLessonInput) error
	DeleteLesson(ctx context.Context, id primitive.ObjectID) error
	AttachPackage(ctx context.Context, modules []primitive.ObjectID, packageId primitive.ObjectID) error
}

type LessonContent interface {
	GetByLessons(ctx context.Context, lessonIds []primitive.ObjectID) ([]domain.LessonContent, error)
	GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.LessonContent, error)
	Update(ctx context.Context, lessonId primitive.ObjectID, content string) error
}

type UpdatePackageInput struct {
	ID          primitive.ObjectID
	Name        string
	Description string
}

type Packages interface {
	Create(ctx context.Context, pkg domain.Package) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdatePackageInput) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Package, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Package, error)
}

type UpdateOfferInput struct {
	ID          primitive.ObjectID
	Name        string
	Description string
	Price       *domain.Price
	Packages    []primitive.ObjectID
}

type Offers interface {
	Create(ctx context.Context, offer domain.Offer) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdateOfferInput) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Offer, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Offer, error)
	GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Offer, error)
}

type PromoCodes interface {
	GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error)
	GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error)
}

type Orders interface {
	Create(ctx context.Context, order domain.Order) error
	AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Order, error)
}

type Repositories struct {
	Schools        Schools
	Students       Students
	StudentLessons StudentLessons
	Courses        Courses
	Modules        Modules
	Packages       Packages
	LessonContent  LessonContent
	Offers         Offers
	PromoCodes     PromoCodes
	Orders         Orders
	Admins         Admins
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Schools:        NewSchoolsRepo(db),
		Students:       NewStudentsRepo(db),
		StudentLessons: NewStudentLessonsRepo(db),
		Courses:        NewCoursesRepo(db),
		Modules:        NewModulesRepo(db),
		LessonContent:  NewLessonContentRepo(db),
		Offers:         NewOffersRepo(db),
		PromoCodes:     NewPromocodeRepo(db),
		Orders:         NewOrdersRepo(db),
		Admins:         NewAdminsRepo(db),
		Packages:       NewPackagesRepo(db),
	}
}
