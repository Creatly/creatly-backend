package service

import (
	"context"

	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/cache"
)

type SchoolsService struct {
	repo  repository.Schools
	cache cache.Cache
	ttl   int64
}

func NewSchoolsService(repo repository.Schools, cache cache.Cache, ttl int64) *SchoolsService {
	return &SchoolsService{repo: repo, cache: cache, ttl: ttl}
}

func (s *SchoolsService) GetByDomain(ctx context.Context, domainName string) (domain.School, error) {
	if value, err := s.cache.Get(domainName); err == nil {
		return value.(domain.School), nil
	}

	school, err := s.repo.GetByDomain(ctx, domainName)
	if err != nil {
		return domain.School{}, err
	}

	s.cache.Set(domainName, school, s.ttl)
	return school, nil
}

func (s *SchoolsService) UpdateSettings(ctx context.Context, inp UpdateSchoolSettingsInput) error {
	return s.repo.UpdateSettings(ctx, repository.UpdateSchoolSettingsInput{
		SchoolID:    inp.SchoolID,
		Color:       inp.Color,
		Domains:     inp.Domains,
		Email:       inp.Email,
		ContactData: inp.ContactData,
		Pages:       inp.Pages,
	})
}
