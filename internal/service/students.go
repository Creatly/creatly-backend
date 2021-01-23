package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/auth"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type StudentsService struct {
	repo         repository.Students
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	emailService Emails

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewStudentsService(repo repository.Students, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	emailService Emails, accessTTL, refreshTTL time.Duration) *StudentsService {
	return &StudentsService{
		repo:            repo,
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

func (s *StudentsService) SignIn(ctx context.Context, input StudentSignInInput) (Tokens, error) {
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
