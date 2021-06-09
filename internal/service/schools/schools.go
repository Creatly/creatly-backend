package schools

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"github.com/zhashkevych/creatly-backend/pkg/cache"
)

type (
	UpdateSchoolSettingsInput struct {
		Color       string
		Domains     []string
		Email       string
		ContactInfo *domain.ContactInfo
		Pages       *domain.Pages
	}

	Schools interface {
		Create(ctx context.Context, name string) (primitive.ObjectID, error)
		GetByDomain(ctx context.Context, domainName string) (domain.School, error)
		GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error)
		UpdateSettings(ctx context.Context, inp repository.UpdateSchoolSettingsInput) error
	}

	Cache interface {
		Set(key, value interface{}, ttl int64) error
		Get(key interface{}) (interface{}, error)
	}

	SchoolsService struct {
		repo  Schools
		cache Cache
		ttl   int64
	}
)

func NewSchoolsService(repo Schools, cache cache.Cache, ttl int64) *SchoolsService {
	return &SchoolsService{repo: repo, cache: cache, ttl: ttl}
}

func (s *SchoolsService) Create(ctx context.Context, name string) (primitive.ObjectID, error) {
	return s.repo.Create(ctx, name)
}

func (s *SchoolsService) GetByDomain(ctx context.Context, domainName string) (domain.School, error) {
	if value, err := s.cache.Get(domainName); err == nil {
		return value.(domain.School), nil
	}

	school, err := s.repo.GetByDomain(ctx, domainName)
	if err != nil {
		return domain.School{}, err
	}

	err = s.cache.Set(domainName, school, s.ttl)

	return school, err
}

func (s *SchoolsService) UpdateSettings(ctx context.Context, schoolId primitive.ObjectID, inp UpdateSchoolSettingsInput) error {
	return s.repo.UpdateSettings(ctx, repository.UpdateSchoolSettingsInput{
		SchoolID:    schoolId,
		Color:       inp.Color,
		Domains:     inp.Domains,
		Email:       inp.Email,
		ContactInfo: inp.ContactInfo,
		Pages:       inp.Pages,
	})
}
