package service

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/pkg/cache"
	"github.com/zhashkevych/courses-backend/pkg/logger"
)

type SchoolsService struct {
	repo  repository.Schools
	cache cache.Cache
}

func NewSchoolsService(repo repository.Schools, cache cache.Cache) *SchoolsService {
	return &SchoolsService{repo: repo, cache: cache}
}

func (s *SchoolsService) GetByDomain(ctx context.Context, domainName string) (domain.School, error) {
	if value, err := s.cache.Get(domainName); err == nil {
		return value.(domain.School), nil
	}

	logger.Info(domainName)

	school, err := s.repo.GetByDomain(ctx, domainName)
	if err != nil {
		return domain.School{}, err
	}

	s.cache.Set(domainName, school)
	return school, nil
}