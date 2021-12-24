package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	Verify(ctx context.Context, userID primitive.ObjectID, code string) error
	SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error
	AttachSchool(ctx context.Context, userID, schoolID primitive.ObjectID) error
}

type Schools interface {
	Create(ctx context.Context, name string) (primitive.ObjectID, error)
	GetByDomain(ctx context.Context, domainName string) (domain.School, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error)
	UpdateSettings(ctx context.Context, id primitive.ObjectID, inp domain.UpdateSchoolSettingsInput) error
	SetFondyCredentials(ctx context.Context, id primitive.ObjectID, fondy domain.Fondy) error
}

type Students interface {
	Create(ctx context.Context, student *domain.Student) error
	Update(ctx context.Context, inp domain.UpdateStudentInput) error
	Delete(ctx context.Context, schoolId, studentId primitive.ObjectID) error
	GetByCredentials(ctx context.Context, schoolId primitive.ObjectID, email, password string) (domain.Student, error)
	GetByRefreshToken(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (domain.Student, error)
	GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.Student, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID, query domain.GetStudentsQuery) ([]domain.Student, int64, error)
	SetSession(ctx context.Context, studentId primitive.ObjectID, session domain.Session) error
	GiveAccessToModule(ctx context.Context, studentId, moduleId primitive.ObjectID) error
	AttachOffer(ctx context.Context, studentId, offerId primitive.ObjectID, moduleIds []primitive.ObjectID) error
	DetachOffer(ctx context.Context, studentId, offerId primitive.ObjectID, moduleIds []primitive.ObjectID) error
	Verify(ctx context.Context, code string) (domain.Student, error)
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
	GetPublishedById(ctx context.Context, moduleID primitive.ObjectID) (domain.Module, error)
	GetById(ctx context.Context, moduleID primitive.ObjectID) (domain.Module, error)
	GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error)
	Update(ctx context.Context, inp UpdateModuleInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	DeleteByCourse(ctx context.Context, schoolId, courseId primitive.ObjectID) error
	AddLesson(ctx context.Context, schoolId, id primitive.ObjectID, lesson domain.Lesson) error
	GetByLesson(ctx context.Context, lessonID primitive.ObjectID) (domain.Module, error)
	UpdateLesson(ctx context.Context, inp UpdateLessonInput) error
	DeleteLesson(ctx context.Context, schoolId, id primitive.ObjectID) error
	DetachPackageFromAll(ctx context.Context, schoolId, packageId primitive.ObjectID) error
	AttachPackage(ctx context.Context, schoolId, packageId primitive.ObjectID, modules []primitive.ObjectID) error
	AttachSurvey(ctx context.Context, schoolId, id primitive.ObjectID, survey domain.Survey) error
	DetachSurvey(ctx context.Context, schoolId, id primitive.ObjectID) error
}

type LessonContent interface {
	GetByLessons(ctx context.Context, lessonIds []primitive.ObjectID) ([]domain.LessonContent, error)
	GetByLesson(ctx context.Context, lessonID primitive.ObjectID) (domain.LessonContent, error)
	Update(ctx context.Context, schoolID, lessonID primitive.ObjectID, content string) error
	DeleteContent(ctx context.Context, schoolID primitive.ObjectID, lessonIds []primitive.ObjectID) error
}

type UpdatePackageInput struct {
	ID       primitive.ObjectID
	SchoolID primitive.ObjectID
	Name     string
}

type Packages interface {
	Create(ctx context.Context, pkg domain.Package) (primitive.ObjectID, error)
	Update(ctx context.Context, inp UpdatePackageInput) error
	Delete(ctx context.Context, schoolID, id primitive.ObjectID) error
	GetByCourse(ctx context.Context, courseID primitive.ObjectID) ([]domain.Package, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Package, error)
	GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.Package, error)
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
	GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.Offer, error)
}

type PromoCodes interface {
	Create(ctx context.Context, promocode domain.PromoCode) (primitive.ObjectID, error)
	Update(ctx context.Context, inp domain.UpdatePromoCodeInput) error
	Delete(ctx context.Context, schoolId, id primitive.ObjectID) error
	GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error)
	GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.PromoCode, error)
}

type Orders interface {
	Create(ctx context.Context, order domain.Order) error
	AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error)
	GetBySchool(ctx context.Context, schoolId primitive.ObjectID, pagination domain.GetOrdersQuery) ([]domain.Order, int64, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Order, error)
	SetStatus(ctx context.Context, id primitive.ObjectID, status string) error
}

type Files interface {
	Create(ctx context.Context, file domain.File) (primitive.ObjectID, error)
	UpdateStatus(ctx context.Context, fileName string, status domain.FileStatus) error
	GetForUploading(ctx context.Context) (domain.File, error)
	UpdateStatusAndSetURL(ctx context.Context, id primitive.ObjectID, url string) error
	GetByID(ctx context.Context, id, schoolID primitive.ObjectID) (domain.File, error)
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

func filterDateQueries(dateFrom, dateTo, fieldName string, filter bson.M) error {
	if dateFrom != "" && dateTo != "" {
		dateFrom, err := time.Parse(time.RFC3339, dateFrom)
		if err != nil {
			return err
		}

		dateTo, err := time.Parse(time.RFC3339, dateTo)
		if err != nil {
			return err
		}

		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			"$and": []bson.M{
				{fieldName: bson.M{"$gte": dateFrom}},
				{fieldName: bson.M{"$lte": dateTo}},
			},
		})
	}

	if dateFrom != "" && dateTo == "" {
		dateFrom, err := time.Parse(time.RFC3339, dateFrom)
		if err != nil {
			return err
		}

		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			fieldName: bson.M{"$gte": dateFrom},
		})
	}

	if dateFrom == "" && dateTo != "" {
		dateTo, err := time.Parse(time.RFC3339, dateTo)
		if err != nil {
			return err
		}

		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			fieldName: bson.M{"$lte": dateTo},
		})
	}

	return nil
}
