package service

import (
	"context"
	"github.com/zhashkevych/creatly-backend/internal/repository"
)

type SurveysService struct {
	modulesRepo repository.Modules
}

func NewSurveysService(modulesRepo repository.Modules) *SurveysService {
	return &SurveysService{modulesRepo: modulesRepo}
}

func (s *SurveysService) Create(ctx context.Context, inp CreateSurveyInput) error {
	return s.modulesRepo.AttachSurvey(ctx, inp.SchoolID, inp.ModuleID, inp.Survey)
}
