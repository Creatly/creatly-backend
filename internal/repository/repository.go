package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Verify(ctx context.Context, userId primitive.ObjectID, code string) error
	SetSession(ctx context.Context, userId primitive.ObjectID, session domain.Session) error
	AttachSchool(ctx context.Context, userId, schoolId primitive.ObjectID) error
}

type UpdateSchoolSettingsInput struct {
	SchoolID            primitive.ObjectID
	Color               *string
	Domains             []string
	Email               *string
	ContactInfo         *domain.ContactInfo
	Pages               *domain.Pages
	ShowPaymentImages   *bool
	GoogleAnalyticsCode *string
	LogoURL             *string
}

type Schools interface {
	Create(ctx context.Context, name string) (primitive.ObjectID, error)
	GetByDomain(ctx context.Context, domainName string) (domain.School, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error)
	UpdateSettings(ctx context.Context, inp UpdateSchoolSettingsInput) error
	SetFondyCredentials(ctx context.Context, id primitive.ObjectID, fondy domain.Fondy) error
}

type Students interface {
	Create(ctx context.Context, student domain.Student) error
	GetByCredentials(ctx context.Context, schoolId primitive.ObjectID, email, password string) (domain.Student, error)
	GetByRefreshToken(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (domain.Student, error)
	GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.Student, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID, pagination *domain.PaginationQuery) ([]domain.Student, int64, error)
	SetSession(ctx context.Context, studentId primitive.ObjectID, session domain.Session) error
	GiveAccessToModule(ctx context.Context, studentId, moduleId primitive.ObjectID) error
	AttachOffer(ctx context.Context, studentId, offerId primitive.ObjectID, moduleIds []primitive.ObjectID) error
	DetachOffer(ctx context.Context, studentId, offerId primitive.ObjectID, moduleIds []primitive.ObjectID) error
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
	SchoolID    primitive.ObjectID
	Name        *string
	ImageURL    *string
	Description *string
	Color       *string
	Published   *bool
}

type Courses interface {
	Create(ctx context.Context, schoolId primitive.ObjectID, course domain.Course) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdateCourseInput) error
	Delete(ctx context.Context, schoolId, courseId primitive.ObjectID) error
}

type UpdateModuleInput struct {
	ID        primitive.ObjectID
	SchoolID  primitive.ObjectID
	Name      string
	Position  *uint
	Published *bool
}

type UpdateLessonInput struct {
	ID        primitive.ObjectID
	SchoolID  primitive.ObjectID
	Name      string
	Position  *uint
	Published *bool
}

type Modules interface {
	Create(ctx context.Context, module domain.Module) (primitive.ObjectID, error)
	GetPublishedByCourseId(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
	GetByCourseId(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
	GetPublishedById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error)
	Update(ctx context.Context, inp UpdateModuleInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	DeleteByCourse(ctx context.Context, schoolId, courseId primitive.ObjectID) error
	AddLesson(ctx context.Context, schoolId, id primitive.ObjectID, lesson domain.Lesson) error
	GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.Module, error)
	UpdateLesson(ctx context.Context, inp UpdateLessonInput) error
	DeleteLesson(ctx context.Context, schoolId, id primitive.ObjectID) error
	AttachPackage(ctx context.Context, schoolId, packageId primitive.ObjectID, modules []primitive.ObjectID) error
	AttachSurvey(ctx context.Context, schoolId, id primitive.ObjectID, survey domain.Survey) error
	DetachSurvey(ctx context.Context, schoolId, id primitive.ObjectID) error
}

type LessonContent interface {
	GetByLessons(ctx context.Context, lessonIds []primitive.ObjectID) ([]domain.LessonContent, error)
	GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.LessonContent, error)
	Update(ctx context.Context, schoolId, lessonId primitive.ObjectID, content string) error
	DeleteContent(ctx context.Context, schoolId primitive.ObjectID, lessonIds []primitive.ObjectID) error
}

type UpdatePackageInput struct {
	ID          primitive.ObjectID
	SchoolID    primitive.ObjectID
	Name        string
	Description string
}

type Packages interface {
	Create(ctx context.Context, pkg domain.Package) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdatePackageInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Package, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Package, error)
}

type UpdateOfferInput struct {
	ID            primitive.ObjectID
	SchoolID      primitive.ObjectID
	Name          string
	Description   string
	Benefits      []string
	Price         *domain.Price
	Packages      []primitive.ObjectID
	PaymentMethod *domain.PaymentMethod
}

type Offers interface {
	Create(ctx context.Context, offer domain.Offer) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdateOfferInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Offer, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Offer, error)
	GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Offer, error)
}

type UpdatePromoCodeInput struct {
	ID                 primitive.ObjectID
	SchoolID           primitive.ObjectID
	Code               string
	DiscountPercentage int
	ExpiresAt          time.Time
	OfferIDs           []primitive.ObjectID
}

type PromoCodes interface {
	Create(ctx context.Context, promocode domain.PromoCode) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdatePromoCodeInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error)
	GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.PromoCode, error)
}

type Orders interface {
	Create(ctx context.Context, order domain.Order) error
	AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID, pagination *domain.PaginationQuery) ([]domain.Order, int64, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Order, error)
	SetStatus(ctx context.Context, id primitive.ObjectID, status string) error
}

type Files interface {
	Create(ctx context.Context, file domain.File) (primitive.ObjectID, error)
	UpdateStatus(ctx context.Context, fileName string, status domain.FileStatus) error
	GetForUploading(ctx context.Context) (domain.File, error)
	UpdateStatusAndSetURL(ctx context.Context, id primitive.ObjectID, url string) error
	GetByID(ctx context.Context, id, schoolId primitive.ObjectID) (domain.File, error)
}

type SurveyResults interface {
	Save(ctx context.Context, results domain.SurveyResult) error
	GetAllByModule(ctx context.Context, moduleId primitive.ObjectID, pagination *domain.PaginationQuery) ([]domain.SurveyResult, int64, error)
	GetByStudent(ctx context.Context, moduleId, studentId primitive.ObjectID) (domain.SurveyResult, error)
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
	Users          Users
	Files          Files
	SurveyResults  SurveyResults
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
		Users:          NewUsersRepo(db),
		Files:          NewFilesRepo(db),
		SurveyResults:  NewSurveyResultsRepo(db),
	}
}

func getPaginationOpts(pagination *domain.PaginationQuery) *options.FindOptions {
	var opts *options.FindOptions
	if pagination != nil {
		opts = &options.FindOptions{
			Skip:  pagination.GetSkip(),
			Limit: pagination.GetLimit(),
		}
	}

	return opts
}
