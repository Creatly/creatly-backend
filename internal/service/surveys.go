package service

import (
	"context"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SurveysService struct {
	modulesRepo repository.Modules
}

func NewSurveysService(modulesRepo repository.Modules) *SurveysService {
	return &SurveysService{modulesRepo: modulesRepo}
}

func (s *SurveysService) Create(ctx context.Context, inp CreateSurveyInput) error {
	for i := range inp.Survey.Questions {
		inp.Survey.Questions[i].ID = primitive.NewObjectID()
	}

	return s.modulesRepo.AttachSurvey(ctx, inp.SchoolID, inp.ModuleID, inp.Survey)
}
