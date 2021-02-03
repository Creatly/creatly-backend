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

type AdminsService struct {
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager
	repo         repository.Admins

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAdminsService(hasher hash.PasswordHasher, tokenManager auth.TokenManager, repo repository.Admins, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *AdminsService {
	return &AdminsService{hasher: hasher, tokenManager: tokenManager, repo: repo, accessTokenTTL: accessTokenTTL, refreshTokenTTL: refreshTokenTTL}
}

func (s *AdminsService) SignIn(ctx context.Context, input SignInInput) (Tokens, error) {
	//student, err := s.repo.GetByCredentials(ctx, input.SchoolID, input.Email, s.hasher.Hash(input.Password))
	student, err := s.repo.GetByCredentials(ctx, input.SchoolID, input.Email, input.Password) // TODO implement password hashing
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, student.ID)
}

func (s *AdminsService) RefreshTokens(ctx context.Context, schoolId primitive.ObjectID, refreshToken string) (Tokens, error) {
	student, err := s.repo.GetByRefreshToken(ctx, schoolId, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, student.ID)
}

func (s *AdminsService) createSession(ctx context.Context, adminId primitive.ObjectID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(adminId.Hex(), s.accessTokenTTL)
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

	err = s.repo.SetSession(ctx, adminId, session)
	return res, err
}
