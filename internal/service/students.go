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
	repo           repository.Students
	coursesService Courses
	hasher         hash.PasswordHasher
	tokenManager   auth.TokenManager
	emailService   Emails

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewStudentsService(repo repository.Students, coursesService Courses, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	emailService Emails, accessTTL, refreshTTL time.Duration) *StudentsService {
	return &StudentsService{
		repo:            repo,
		coursesService:  coursesService,
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
		Name:           input.Name,
		Password:       s.hasher.Hash(input.Password),
		Email:          input.Email,
		RegisteredAt:   time.Now(),
		LastVisitAt:    time.Now(),
		SchoolID:       input.SchoolID,
		RegisterSource: input.RegisterSource,
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
		RegisterSource:   student.RegisterSource,
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

func (s *StudentsService) GetStudentModuleWithLessons(ctx context.Context, schoolId, studentId, moduleId primitive.ObjectID) ([]domain.Lesson, error) {
	// Get module with lessons content, check if it is available for student
	module, err := s.coursesService.GetModuleWithContent(ctx, moduleId)
	if err != nil {
		return nil, err
	}

	logger.Info(module)

	student, err := s.repo.GetById(ctx, studentId)
	if err != nil {
		return nil, err
	}

	logger.Info(student)

	if student.IsModuleAvailable(module) {
		logger.Info("Module is available")
		return module.Lessons, nil
	}

	logger.Info("Ooops")

	// Find module offers
	offers, err := s.coursesService.GetPackageOffers(ctx, schoolId, module.PackageID)
	if err != nil {
		return nil, err
	}

	if len(offers) != 0 {
		return nil, ErrModuleIsNotAvailable
	}

	// If module has no offers - it's free and available to everyone
	go func() {
		if err := s.repo.GiveAccessToModules(ctx, studentId, []primitive.ObjectID{moduleId}); err != nil {
			logger.Error(err)
		}
	}()

	return module.Lessons, nil
}

func (s *StudentsService) GiveAccessToModules(ctx context.Context, studentId primitive.ObjectID, moduleIds []primitive.ObjectID) error {
	return s.repo.GiveAccessToModules(ctx, studentId, moduleIds)
}

func (s *StudentsService) GiveAccessToPackages(ctx context.Context, studentId primitive.ObjectID, packageIds []primitive.ObjectID) error {
	modules, err := s.coursesService.GetPackagesModules(ctx, packageIds)
	if err != nil {
		return err
	}

	ids := make([]primitive.ObjectID, len(modules))
	for i := range modules {
		ids[i] = modules[i].ID
	}

	return s.repo.GiveAccessToModules(ctx, studentId, ids)
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
