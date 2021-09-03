package service

import (
	"context"
	"errors"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"github.com/zhashkevych/creatly-backend/pkg/auth"
	"github.com/zhashkevych/creatly-backend/pkg/hash"
	"github.com/zhashkevych/creatly-backend/pkg/logger"
	"github.com/zhashkevych/creatly-backend/pkg/otp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentsService struct {
	repo         repository.Students
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	otpGenerator otp.Generator

	modulesService        Modules
	offersService         Offers
	emailService          Emails
	lessonsService        Lessons
	studentLessonsService StudentLessons

	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
	verificationCodeLength int
}

func NewStudentsService(repo repository.Students, modulesService Modules, offersService Offers, lessonsService Lessons, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	emailService Emails, studentLessonsService StudentLessons, accessTTL, refreshTTL time.Duration, otpGenerator otp.Generator, verificationCodeLength int) *StudentsService {
	return &StudentsService{
		repo:                   repo,
		modulesService:         modulesService,
		offersService:          offersService,
		hasher:                 hasher,
		emailService:           emailService,
		lessonsService:         lessonsService,
		studentLessonsService:  studentLessonsService,
		tokenManager:           tokenManager,
		accessTokenTTL:         accessTTL,
		refreshTokenTTL:        refreshTTL,
		otpGenerator:           otpGenerator,
		verificationCodeLength: verificationCodeLength,
	}
}

func (s *StudentsService) SignUp(ctx context.Context, input StudentSignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	student := domain.Student{
		Name:         input.Name,
		Password:     passwordHash,
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
		SchoolID:     input.SchoolID,
	}

	if input.Verified {
		student.Verification.Verified = true

		go s.addStudentToList(context.Background(), student)

		return s.repo.Create(ctx, student)
	}

	// it's possible to use OTP apps (Google Authenticator, Authy) compatibility mode here, in the future
	verificationCode := s.otpGenerator.RandomSecret(s.verificationCodeLength)
	student.Verification.Code = verificationCode

	if err := s.repo.Create(ctx, student); err != nil {
		return err
	}

	// TODO: If it fails, what then?
	return s.emailService.SendStudentVerificationEmail(VerificationEmailInput{
		Email:            student.Email,
		Name:             student.Name,
		VerificationCode: verificationCode,
		Domain:           input.SchoolDomain,
	})
}

func (s *StudentsService) SignIn(ctx context.Context, input SchoolSignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	student, err := s.repo.GetByCredentials(ctx, input.SchoolID, input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return Tokens{}, err
		}

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
	student, err := s.repo.Verify(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			return domain.ErrVerificationCodeInvalid
		}

		return err
	}

	logger.Info(student)

	go s.addStudentToList(context.Background(), student)

	return nil
}

func (s *StudentsService) GetModuleContent(ctx context.Context, schoolId, studentId, moduleId primitive.ObjectID) (domain.ModuleContent, error) {
	// Get module with lessons content, check if it is available for student
	module, err := s.modulesService.GetWithContent(ctx, moduleId)
	if err != nil {
		return domain.ModuleContent{}, err
	}

	student, err := s.repo.GetById(ctx, schoolId, studentId)
	if err != nil {
		return domain.ModuleContent{}, err
	}

	if student.IsModuleAvailable(module) {
		return domain.ModuleContent{
			Lessons: module.Lessons,
			Survey:  module.Survey,
		}, nil
	}

	// Find module offers
	offers, err := s.offersService.GetByModule(ctx, schoolId, module.ID)
	if err != nil {
		return domain.ModuleContent{}, err
	}

	if len(offers) != 0 {
		return domain.ModuleContent{}, domain.ErrModuleIsNotAvailable
	}

	// If module has no offers - it's free and available to everyone
	if err := s.repo.GiveAccessToModule(ctx, studentId, moduleId); err != nil {
		return domain.ModuleContent{}, err
	}

	return domain.ModuleContent{
		Lessons: module.Lessons,
		Survey:  module.Survey,
	}, nil
}

func (s *StudentsService) GetLesson(ctx context.Context, studentId, lessonId primitive.ObjectID) (domain.Lesson, error) {
	if err := s.isLessonAvailable(ctx, studentId, lessonId); err != nil {
		return domain.Lesson{}, err
	}

	lesson, err := s.lessonsService.GetById(ctx, lessonId)
	if err != nil {
		return domain.Lesson{}, err
	}

	if err := s.studentLessonsService.SetLastOpened(ctx, studentId, lessonId); err != nil {
		return domain.Lesson{}, err
	}

	return lesson, nil
}

func (s *StudentsService) SetLessonFinished(ctx context.Context, studentId, lessonId primitive.ObjectID) error {
	if err := s.isLessonAvailable(ctx, studentId, lessonId); err != nil {
		return err
	}

	return s.studentLessonsService.AddFinished(ctx, studentId, lessonId)
}

func (s *StudentsService) GiveAccessToOffer(ctx context.Context, studentId primitive.ObjectID, offer domain.Offer) error {
	modules, err := s.modulesService.GetByPackages(ctx, offer.PackageIDs)
	if err != nil {
		return err
	}

	moduleIds := make([]primitive.ObjectID, len(modules))

	for i := range modules {
		moduleIds[i] = modules[i].ID
	}

	return s.repo.AttachOffer(ctx, studentId, offer.ID, moduleIds)
}

func (s *StudentsService) RemoveAccessToOffer(ctx context.Context, studentId primitive.ObjectID, offer domain.Offer) error {
	modules, err := s.modulesService.GetByPackages(ctx, offer.PackageIDs)
	if err != nil {
		return err
	}

	moduleIds := make([]primitive.ObjectID, len(modules))

	for i := range modules {
		moduleIds[i] = modules[i].ID
	}

	return s.repo.DetachOffer(ctx, studentId, offer.ID, moduleIds)
}

func (s *StudentsService) GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.Student, error) {
	return s.repo.GetById(ctx, schoolId, id)
}

func (s *StudentsService) GetBySchool(ctx context.Context, schoolId primitive.ObjectID, pagination *domain.PaginationQuery) ([]domain.Student, int64, error) {
	return s.repo.GetBySchool(ctx, schoolId, pagination)
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

func (s *StudentsService) isLessonAvailable(ctx context.Context, studentId, lessonId primitive.ObjectID) error {
	module, err := s.modulesService.GetByLesson(ctx, lessonId)
	if err != nil {
		return err
	}

	student, err := s.GetById(ctx, module.SchoolID, studentId)
	if err != nil {
		return err
	}

	if !student.IsModuleAvailable(module) {
		return domain.ErrModuleIsNotAvailable
	}

	return nil
}

// TODO refactor.
func (s *StudentsService) addStudentToList(ctx context.Context, student domain.Student) {
	if err := s.emailService.AddStudentToList(ctx, student.Email, student.Name, student.SchoolID); err != nil {
		if err == domain.ErrSendPulseIsNotConnected {
			return
		}

		logger.Errorf("[SENDPULSE] failed to add email to the list: %s", err.Error())
	}
}
