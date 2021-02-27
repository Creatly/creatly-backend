package service

import (
	"context"
	"time"

	"github.com/zhashkevych/courses-backend/pkg/otp"

	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/auth"
	"github.com/zhashkevych/courses-backend/pkg/cache"
	"github.com/zhashkevych/courses-backend/pkg/email"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"github.com/zhashkevych/courses-backend/pkg/payment"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

// TODO handle "not found" errors

type UpdateSchoolSettingsInput struct {
	SchoolID    primitive.ObjectID
	Color       string
	Domain      string
	Email       string
	ContactData string
	Pages       *domain.Pages
}

type Schools interface {
	GetByDomain(ctx context.Context, domainName string) (domain.School, error)
	UpdateSettings(ctx context.Context, input UpdateSchoolSettingsInput) error
}

type StudentSignUpInput struct {
	Name     string
	Email    string
	Password string
	SchoolID primitive.ObjectID
}

type SignInInput struct {
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
	SignIn(ctx context.Context, input SignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (Tokens, error)
	Verify(ctx context.Context, hash string) error
	GetModuleLessons(ctx context.Context, schoolId, studentId, moduleId primitive.ObjectID) ([]domain.Lesson, error)
	GetLesson(ctx context.Context, studentId, lessonId primitive.ObjectID) (domain.Lesson, error)
	SetLessonFinished(ctx context.Context, studentId, lessonId primitive.ObjectID) error
	GiveAccessToPackages(ctx context.Context, studentId primitive.ObjectID, packageIds []primitive.ObjectID) error
	GetAvailableCourses(ctx context.Context, school domain.School, studentId primitive.ObjectID) ([]domain.Course, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Student, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Student, error)
}

type StudentLessons interface {
	AddFinished(ctx context.Context, studentId, lessonId primitive.ObjectID) error
	SetLastOpened(ctx context.Context, studentId, lessonId primitive.ObjectID) error
}

type Admins interface {
	SignIn(ctx context.Context, input SignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (Tokens, error)
	GetCourses(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Course, error)
	GetCourseById(ctx context.Context, schoolId, courseId primitive.ObjectID) (domain.Course, error)
}

type AddToListInput struct {
	Email            string
	Name             string
	VerificationCode string
}

type Emails interface {
	AddToList(AddToListInput) error
}

type UpdateCourseInput struct {
	CourseID    string
	Name        string
	Code        string
	Description string
	Color       string
	Published   *bool
}

type Courses interface {
	Create(ctx context.Context, schoolId primitive.ObjectID, name string) (primitive.ObjectID, error)
	Update(ctx context.Context, schoolId primitive.ObjectID, inp UpdateCourseInput) error
}

type PromoCodes interface {
	GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error)
	GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error)
}

type CreateOfferInput struct {
	Name        string
	Description string
	SchoolID    primitive.ObjectID
	Price       domain.Price
}

type UpdateOfferInput struct {
	ID          string
	Name        string
	Description string
	Price       *domain.Price
	Packages    []string
}

type Offers interface {
	Create(ctx context.Context, inp CreateOfferInput) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdateOfferInput) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Offer, error)
	GetByModule(ctx context.Context, schoolId, moduleId primitive.ObjectID) ([]domain.Offer, error)
	GetByPackage(ctx context.Context, schoolId, packageId primitive.ObjectID) ([]domain.Offer, error)
	GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Offer, error)
	GetAll(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Offer, error)
}

type CreateModuleInput struct {
	CourseID string
	Name     string
	Position uint
}

type UpdateModuleInput struct {
	ID        string
	Name      string
	Position  *uint
	Published *bool
}

type Modules interface {
	Create(ctx context.Context, inp CreateModuleInput) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdateModuleInput) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
	GetById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error)
	GetWithContent(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.Module, error)
}

type AddLessonInput struct {
	ModuleID string
	Name     string
	Position uint
}

type UpdateLessonInput struct {
	LessonID  string
	Name      string
	Content   string
	Position  *uint
	Published *bool
}

type Lessons interface {
	Create(ctx context.Context, inp AddLessonInput) (primitive.ObjectID, error)
	GetById(ctx context.Context, lessonId primitive.ObjectID) (domain.Lesson, error)
	Update(ctx context.Context, inp UpdateLessonInput) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type CreatePackageInput struct {
	CourseID    string
	Name        string
	Description string
}

type UpdatePackageInput struct {
	ID          string
	Name        string
	Description string
	Modules     []string
}

type Packages interface {
	Create(ctx context.Context, inp CreatePackageInput) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdatePackageInput) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Package, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Package, error)
}

type Orders interface {
	Create(ctx context.Context, studentId, offerId, promocodeId primitive.ObjectID) (string, error)
	AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Order, error)
}

type Payments interface {
	ProcessTransaction(ctx context.Context, callback interface{}) error
}

type Services struct {
	Schools        Schools
	Students       Students
	StudentLessons StudentLessons
	Courses        Courses
	PromoCodes     PromoCodes
	Offers         Offers
	Packages       Packages
	Modules        Modules
	Lessons        Lessons
	Payments       Payments
	Orders         Orders
	Admins         Admins
}

type Deps struct {
	Repos                  *repository.Repositories
	Cache                  cache.Cache
	Hasher                 hash.PasswordHasher
	TokenManager           auth.TokenManager
	EmailProvider          email.Provider
	EmailListId            string
	PaymentProvider        payment.Provider
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	PaymentCallbackURL     string
	PaymentResponseURL     string
	CacheTTL               int64
	OtpGenerator           otp.Generator
	VerificationCodeLength int
	FrontendURL            string
}

func NewServices(deps Deps) *Services {
	emailsService := NewEmailsService(deps.EmailProvider, deps.EmailListId, deps.FrontendURL)
	coursesService := NewCoursesService(deps.Repos.Courses)
	modulesService := NewModulesService(deps.Repos.Modules, deps.Repos.LessonContent)
	packagesService := NewPackagesService(deps.Repos.Packages, deps.Repos.Modules)
	offersService := NewOffersService(deps.Repos.Offers, modulesService, packagesService)
	promoCodesService := NewPromoCodeService(deps.Repos.PromoCodes)
	lessonsService := NewLessonsService(deps.Repos.Modules, deps.Repos.LessonContent)
	studentLessonsService := NewStudentLessonsService(deps.Repos.StudentLessons)
	studentsService := NewStudentsService(deps.Repos.Students, modulesService, offersService, lessonsService, deps.Hasher,
		deps.TokenManager, emailsService, studentLessonsService, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.OtpGenerator, deps.VerificationCodeLength)
	ordersService := NewOrdersService(deps.Repos.Orders, offersService, promoCodesService, studentsService, deps.PaymentProvider, deps.PaymentCallbackURL, deps.PaymentResponseURL)

	return &Services{
		Schools:        NewSchoolsService(deps.Repos.Schools, deps.Cache, deps.CacheTTL),
		Students:       studentsService,
		StudentLessons: studentLessonsService,
		Courses:        coursesService,
		PromoCodes:     promoCodesService,
		Offers:         offersService,
		Modules:        modulesService,
		Payments:       NewPaymentsService(deps.PaymentProvider, ordersService, offersService, studentsService),
		Orders:         ordersService,
		Admins:         NewAdminsService(deps.Hasher, deps.TokenManager, deps.Repos.Admins, deps.Repos.Schools, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Packages:       packagesService,
		Lessons:        lessonsService,
	}
}
