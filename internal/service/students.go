package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/auth"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"github.com/zhashkevych/courses-backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type StudentsService struct {
	repo         repository.Students
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	modulesService Modules
	offersService  Offers
	emailService   Emails

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewStudentsService(repo repository.Students, modulesService Modules, offersService Offers, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	emailService Emails, accessTTL, refreshTTL time.Duration) *StudentsService {
	return &StudentsService{
		repo:            repo,
		modulesService:  modulesService,
		offersService:   offersService,
		hasher:          hasher,
		emailService:    emailService,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *StudentsService) SignUp(ctx context.Context, input StudentSignUpInput) error {
	verificationCode := primitive.NewObjectID()
	student := domain.Student{
		Name:         input.Name,
		Password:     s.hasher.Hash(input.Password),
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
		SchoolID:     input.SchoolID,
		Verification: domain.Verification{
			Code: verificationCode,
		},
	}

	if err := s.repo.Create(ctx, student); err != nil {
		return err
	}

	// TODO: If it fails, what then?
	return s.emailService.AddToList(AddToListInput{
		Email:            student.Email,
		Name:             student.Name,
		VerificationCode: verificationCode.Hex(),
	})
}

func (s *StudentsService) SignIn(ctx context.Context, input SignInInput) (Tokens, error) {
	student, err := s.repo.GetByCredentials(ctx, input.SchoolID, input.Email, s.hasher.Hash(input.Password))
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, student.ID)
}

func (s *StudentsService) RefreshTokens(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (Tokens, error) {
	student, err := s.repo.GetByRefreshToken(ctx, schoolId, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, student.ID)
}

func (s *StudentsService) Verify(ctx context.Context, hash string) error {
	return s.repo.Verify(ctx, hash)
}

func (s *StudentsService) GetModuleLessons(ctx context.Context, schoolId, studentId, moduleId primitive.ObjectID) ([]domain.Lesson, error) {
	// Get module with lessons content, check if it is available for student
	module, err := s.modulesService.GetWithContent(ctx, moduleId)
	if err != nil {
		return nil, err
	}

	student, err := s.repo.GetById(ctx, studentId)
	if err != nil {
		return nil, err
	}

	if student.IsModuleAvailable(module) {
		logger.Info("Module is available")
		return module.Lessons, nil
	}

	// Find module offers
	offers, err := s.offersService.GetByModule(ctx, schoolId, module.ID)
	if err != nil {
		return nil, err
	}

	if len(offers) != 0 {
		return nil, ErrModuleIsNotAvailable
	}

	// If module has no offers - it's free and available to everyone
	go func() {
		if err := s.repo.GiveAccessToCourseAndModule(ctx, studentId, module.CourseID, moduleId); err != nil {
			logger.Error(err)
		}
	}()

	return module.Lessons, nil
}

func (s *StudentsService) GiveAccessToPackages(ctx context.Context, studentId primitive.ObjectID, packageIds []primitive.ObjectID) error {
	modules, err := s.modulesService.GetByPackages(ctx, packageIds)
	if err != nil {
		return err
	}

	moduleIds := make([]primitive.ObjectID, len(modules))
	courses := map[primitive.ObjectID]struct{}{}
	for i := range modules {
		moduleIds[i] = modules[i].ID
		courses[modules[i].CourseID] = struct{}{}
	}

	courseIds := make([]primitive.ObjectID, 0)
	for id := range courses {
		courseIds = append(courseIds, id)
	}

	return s.repo.GiveAccessToCoursesAndModules(ctx, studentId, courseIds, moduleIds)
}

func (s *StudentsService) GetAvailableCourses(ctx context.Context, school domain.School, studentId primitive.ObjectID) ([]domain.Course, error) {
	student, err := s.repo.GetById(ctx, studentId)
	if err != nil {
		return nil, err
	}

	courses := make([]domain.Course, 0)
	for _, id := range student.AvailableCourses {
		for _, course := range school.Courses {
			if id == course.ID {
				courses = append(courses, course)
			}
		}
	}

	return courses, nil
}

func (s *StudentsService) createSession(ctx context.Context, studentId primitive.ObjectID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(studentId.Hex(), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, studentId, session)
	return res, err
}
