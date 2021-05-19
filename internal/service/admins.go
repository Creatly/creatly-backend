package service

import (
	"context"
	"errors"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"github.com/zhashkevych/creatly-backend/pkg/auth"
	"github.com/zhashkevych/creatly-backend/pkg/hash"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminsService struct {
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	repo       repository.Admins
	schoolRepo repository.Schools

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAdminsService(hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	repo repository.Admins, schoolRepo repository.Schools, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *AdminsService {
	return &AdminsService{
		hasher:          hasher,
		tokenManager:    tokenManager,
		repo:            repo,
		schoolRepo:      schoolRepo,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AdminsService) SignIn(ctx context.Context, input SignInInput) (Tokens, error) {
	// student, err := s.repo.GetByCredentials(ctx, input.SchoolID, input.Email, s.hasher.Hash(input.Password))
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

func (s *AdminsService) GetCourses(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Course, error) {
	school, err := s.schoolRepo.GetById(ctx, schoolId)
	if err != nil {
		return nil, err
	}

	return school.Courses, nil
}

func (s *AdminsService) GetCourseById(ctx context.Context, schoolId, courseId primitive.ObjectID) (domain.Course, error) {
	school, err := s.schoolRepo.GetById(ctx, schoolId)
	if err != nil {
		return domain.Course{}, err
	}

	var searchedCourse domain.Course
	for _, course := range school.Courses {
		if course.ID == courseId {
			searchedCourse = course
		}
	}

	if searchedCourse.ID.IsZero() {
		return domain.Course{}, errors.New("not found")
	}

	return searchedCourse, nil
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
