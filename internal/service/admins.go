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

	repo        repository.Admins
	schoolRepo  repository.Schools
	studentRepo repository.Students

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAdminsService(hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	repo repository.Admins, schoolRepo repository.Schools, studentRepo repository.Students,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *AdminsService {
	return &AdminsService{
		hasher:          hasher,
		tokenManager:    tokenManager,
		repo:            repo,
		schoolRepo:      schoolRepo,
		studentRepo:     studentRepo,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AdminsService) SignIn(ctx context.Context, input SchoolSignInInput) (Tokens, error) {
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

func (s *AdminsService) GetCourseById(ctx context.Context, schoolID, courseID primitive.ObjectID) (domain.Course, error) {
	school, err := s.schoolRepo.GetById(ctx, schoolID)
	if err != nil {
		return domain.Course{}, err
	}

	var searchedCourse domain.Course

	for _, course := range school.Courses {
		if course.ID == courseID {
			searchedCourse = course
		}
	}

	if searchedCourse.ID.IsZero() {
		return domain.Course{}, errors.New("not found")
	}

	return searchedCourse, nil
}

func (s *AdminsService) CreateStudent(ctx context.Context, inp domain.CreateStudentInput) (domain.Student, error) {
	passwordHash, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return domain.Student{}, err
	}

	student := domain.Student{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     passwordHash,
		RegisteredAt: time.Now(),
		SchoolID:     inp.SchoolID,
		Verification: domain.Verification{Verified: true},
	}
	err = s.studentRepo.Create(ctx, &student)

	return student, err
}

func (s *AdminsService) UpdateStudent(ctx context.Context, inp domain.UpdateStudentInput) error {
	return s.studentRepo.Update(ctx, inp)
}

func (s *AdminsService) DeleteStudent(ctx context.Context, schoolId, studentId primitive.ObjectID) error {
	return s.studentRepo.Delete(ctx, schoolId, studentId)
}

func (s *AdminsService) createSession(ctx context.Context, adminID primitive.ObjectID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(adminID.Hex(), s.accessTokenTTL)
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

	err = s.repo.SetSession(ctx, adminID, session)

	return res, err
}
