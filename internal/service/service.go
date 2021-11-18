package service

import (
	"context"
	"io"
	"time"

	"github.com/zhashkevych/creatly-backend/pkg/dns"

	"github.com/zhashkevych/creatly-backend/internal/config"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"github.com/zhashkevych/creatly-backend/pkg/auth"
	"github.com/zhashkevych/creatly-backend/pkg/cache"
	"github.com/zhashkevych/creatly-backend/pkg/email"
	"github.com/zhashkevych/creatly-backend/pkg/hash"
	"github.com/zhashkevych/creatly-backend/pkg/otp"
	"github.com/zhashkevych/creatly-backend/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

// TODO handle "not found" errors

type UserSignUpInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

// 1. Create School in DB
// 2. Generate Sub Domain

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	Verify(ctx context.Context, userID primitive.ObjectID, hash string) error
	CreateSchool(ctx context.Context, userID primitive.ObjectID, schoolName string) (domain.School, error)
}

type ConnectFondyInput struct {
	SchoolID         primitive.ObjectID
	MerchantID       string
	MerchantPassword string
}

type ConnectSendPulseInput struct {
	SchoolID primitive.ObjectID
	ID       string
	Secret   string
	ListID   string
}

type Schools interface {
	Create(ctx context.Context, name string) (primitive.ObjectID, error)
	GetByDomain(ctx context.Context, domainName string) (domain.School, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error)
	UpdateSettings(ctx context.Context, schoolId primitive.ObjectID, input domain.UpdateSchoolSettingsInput) error
	ConnectFondy(ctx context.Context, input ConnectFondyInput) error
	ConnectSendPulse(ctx context.Context, input ConnectSendPulseInput) error
}

type StudentSignUpInput struct {
	Name         string
	Email        string
	Password     string
	SchoolID     primitive.ObjectID
	SchoolDomain string
	Verified     bool
}

type SchoolSignInInput struct {
	Email    string
	Password string
	SchoolID primitive.ObjectID
}

type Students interface {
	SignUp(ctx context.Context, input StudentSignUpInput) error
	SignIn(ctx context.Context, input SchoolSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (Tokens, error)
	Verify(ctx context.Context, hash string) error
	GetModuleContent(ctx context.Context, schoolId, studentId, moduleId primitive.ObjectID) (domain.ModuleContent, error)
	GetLesson(ctx context.Context, studentId, lessonId primitive.ObjectID) (domain.Lesson, error)
	SetLessonFinished(ctx context.Context, studentId, lessonId primitive.ObjectID) error
	GiveAccessToOffer(ctx context.Context, studentId primitive.ObjectID, offer domain.Offer) error
	RemoveAccessToOffer(ctx context.Context, studentId primitive.ObjectID, offer domain.Offer) error
	GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.Student, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID, query domain.GetStudentsQuery) ([]domain.Student, int64, error)
}

type StudentLessons interface {
	AddFinished(ctx context.Context, studentId, lessonId primitive.ObjectID) error
	SetLastOpened(ctx context.Context, studentId, lessonId primitive.ObjectID) error
}

type Admins interface {
	SignIn(ctx context.Context, input SchoolSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (Tokens, error)
	GetCourses(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Course, error)
	GetCourseById(ctx context.Context, schoolId, courseId primitive.ObjectID) (domain.Course, error)
	CreateStudent(ctx context.Context, inp domain.CreateStudentInput) (domain.Student, error)
	UpdateStudent(ctx context.Context, inp domain.UpdateStudentInput) error
	DeleteStudent(ctx context.Context, schoolId, studentId primitive.ObjectID) error
}

type UploadInput struct {
	File        io.Reader
	Filename    string
	Size        int64
	ContentType string
	SchoolID    primitive.ObjectID
	Type        domain.FileType
}

type Files interface {
	UploadAndSaveFile(ctx context.Context, file domain.File) (string, error)
	Save(ctx context.Context, file domain.File) (primitive.ObjectID, error)
	UpdateStatus(ctx context.Context, fileName string, status domain.FileStatus) error // TODO check schoolID
	GetByID(ctx context.Context, id, schoolId primitive.ObjectID) (domain.File, error)
	InitStorageUploaderWorkers(ctx context.Context)
}

type VerificationEmailInput struct {
	Email            string
	Name             string
	VerificationCode string
	Domain           string
}

type StudentPurchaseSuccessfulEmailInput struct {
	Email      string
	Name       string
	CourseName string
}

type Emails interface {
	SendStudentVerificationEmail(VerificationEmailInput) error
	SendUserVerificationEmail(VerificationEmailInput) error
	SendStudentPurchaseSuccessfulEmail(StudentPurchaseSuccessfulEmailInput) error
	AddStudentToList(ctx context.Context, email, name string, schoolID primitive.ObjectID) error
}

type UpdateCourseInput struct {
	CourseID    string
	SchoolID    string
	Name        *string
	ImageURL    *string
	Description *string
	Color       *string
	Published   *bool
}

type Courses interface {
	Create(ctx context.Context, schoolId primitive.ObjectID, name string) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdateCourseInput) error
	Delete(ctx context.Context, schoolId, courseId primitive.ObjectID) error
}

type CreatePromoCodeInput struct {
	SchoolID           primitive.ObjectID
	Code               string
	DiscountPercentage int
	ExpiresAt          time.Time
	OfferIDs           []primitive.ObjectID
}

type PromoCodes interface {
	Create(ctx context.Context, inp CreatePromoCodeInput) (primitive.ObjectID, error)
	Update(ctx context.Context, inp domain.UpdatePromoCodeInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error)
	GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.PromoCode, error)
}

type CreateOfferInput struct {
	Name          string
	Description   string
	Benefits      []string
	SchoolID      primitive.ObjectID
	Price         domain.Price
	Packages      []string
	PaymentMethod domain.PaymentMethod
}

type UpdateOfferInput struct {
	ID            string
	SchoolID      string
	Name          string
	Description   string
	Benefits      []string
	Price         *domain.Price
	Packages      []string
	PaymentMethod *domain.PaymentMethod
}

func (i UpdateOfferInput) ValidatePayment() error {
	if i.PaymentMethod == nil {
		return nil
	}

	if !i.PaymentMethod.UsesProvider {
		return nil
	}

	return i.PaymentMethod.Validate()
}

type Offers interface {
	Create(ctx context.Context, inp CreateOfferInput) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdateOfferInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Offer, error)
	GetByModule(ctx context.Context, schoolId, moduleId primitive.ObjectID) ([]domain.Offer, error)
	GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Offer, error)
	GetAll(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Offer, error)
	GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.Offer, error)
}

type CreateModuleInput struct {
	SchoolID string
	CourseID string
	Name     string
	Position uint
}

type UpdateModuleInput struct {
	ID        string
	SchoolID  string
	Name      string
	Position  *uint
	Published *bool
}

type Modules interface {
	Create(ctx context.Context, inp CreateModuleInput) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdateModuleInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	DeleteByCourse(ctx context.Context, schoolId, courseId primitive.ObjectID) error
	GetPublishedByCourseId(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
	GetByCourseId(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error)
	GetById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error)
	GetWithContent(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error)
	GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.Module, error)
}

type AddLessonInput struct {
	ModuleID string
	SchoolID string
	Name     string
	Position uint
}

type UpdateLessonInput struct {
	LessonID  string
	SchoolID  string
	Name      string
	Content   string
	Position  *uint
	Published *bool
}

type Lessons interface {
	Create(ctx context.Context, inp AddLessonInput) (primitive.ObjectID, error)
	GetById(ctx context.Context, lessonId primitive.ObjectID) (domain.Lesson, error)
	Update(ctx context.Context, inp UpdateLessonInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	DeleteContent(ctx context.Context, schoolId primitive.ObjectID, lessonIds []primitive.ObjectID) error
}

type CreatePackageInput struct {
	CourseID string
	SchoolID string
	Name     string
	Modules  []string
}

type UpdatePackageInput struct {
	ID       string
	SchoolID string
	Name     string
	Modules  []string
}

type Packages interface {
	Create(ctx context.Context, inp CreatePackageInput) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdatePackageInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Package, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Package, error)
	GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.Package, error)
}

type Orders interface {
	Create(ctx context.Context, studentId, offerId, promocodeId primitive.ObjectID) (primitive.ObjectID, error)
	AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID, query domain.GetOrdersQuery) ([]domain.Order, int64, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Order, error)
	SetStatus(ctx context.Context, id primitive.ObjectID, status string) error
}

type Payments interface {
	GeneratePaymentLink(ctx context.Context, orderId primitive.ObjectID) (string, error)
	ProcessTransaction(ctx context.Context, callback interface{}) error
}

type CreateSurveyInput struct {
	ModuleID primitive.ObjectID
	SchoolID primitive.ObjectID
	Survey   domain.Survey
}

type SaveStudentAnswersInput struct {
	ModuleID  primitive.ObjectID
	StudentID primitive.ObjectID
	SchoolID  primitive.ObjectID
	Answers   []domain.SurveyAnswer
}

type Surveys interface {
	Create(ctx context.Context, inp CreateSurveyInput) error
	Delete(ctx context.Context, schoolId, moduleId primitive.ObjectID) error
	SaveStudentAnswers(ctx context.Context, inp SaveStudentAnswersInput) error
	GetResultsByModule(ctx context.Context, moduleId primitive.ObjectID,
		pagination *domain.PaginationQuery) ([]domain.SurveyResult, int64, error)
	GetStudentResults(ctx context.Context, moduleId, studentId primitive.ObjectID) (domain.SurveyResult, error)
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
	Files          Files
	Users          Users
	Surveys        Surveys
}

type Deps struct {
	Repos                  *repository.Repositories
	Cache                  cache.Cache
	Hasher                 hash.PasswordHasher
	TokenManager           auth.TokenManager
	EmailSender            email.Sender
	EmailConfig            config.EmailConfig
	StorageProvider        storage.Provider
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	FondyCallbackURL       string
	CacheTTL               int64
	OtpGenerator           otp.Generator
	VerificationCodeLength int
	Environment            string
	Domain                 string
	DNS                    dns.DomainManager
}

func NewServices(deps Deps) *Services {
	schoolsService := NewSchoolsService(deps.Repos.Schools, deps.Cache, deps.CacheTTL)
	emailsService := NewEmailsService(deps.EmailSender, deps.EmailConfig, *schoolsService, deps.Cache)
	modulesService := NewModulesService(deps.Repos.Modules, deps.Repos.LessonContent)
	coursesService := NewCoursesService(deps.Repos.Courses, modulesService)
	packagesService := NewPackagesService(deps.Repos.Packages, deps.Repos.Modules)
	offersService := NewOffersService(deps.Repos.Offers, modulesService, packagesService)
	promoCodesService := NewPromoCodeService(deps.Repos.PromoCodes)
	lessonsService := NewLessonsService(deps.Repos.Modules, deps.Repos.LessonContent)
	studentLessonsService := NewStudentLessonsService(deps.Repos.StudentLessons)
	studentsService := NewStudentsService(deps.Repos.Students, modulesService, offersService, lessonsService, deps.Hasher,
		deps.TokenManager, emailsService, studentLessonsService, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.OtpGenerator, deps.VerificationCodeLength)
	ordersService := NewOrdersService(deps.Repos.Orders, offersService, promoCodesService, studentsService)
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, emailsService, schoolsService, deps.DNS,
		deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.OtpGenerator, deps.VerificationCodeLength, deps.Domain)

	return &Services{
		Schools:        schoolsService,
		Students:       studentsService,
		StudentLessons: studentLessonsService,
		Courses:        coursesService,
		PromoCodes:     promoCodesService,
		Offers:         offersService,
		Modules:        modulesService,
		Payments: NewPaymentsService(ordersService, offersService, studentsService, emailsService, schoolsService,
			deps.FondyCallbackURL),
		Orders: ordersService,
		Admins: NewAdminsService(deps.Hasher, deps.TokenManager, deps.Repos.Admins, deps.Repos.Schools, deps.Repos.Students,
			deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Packages: packagesService,
		Lessons:  lessonsService,
		Files:    NewFilesService(deps.Repos.Files, deps.StorageProvider, deps.Environment),
		Users:    usersService,
		Surveys:  NewSurveysService(deps.Repos.Modules, deps.Repos.SurveyResults, deps.Repos.Students),
	}
}
