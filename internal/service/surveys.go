package service

import (
	"context"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SurveysService struct {
	modulesRepo       repository.Modules
	surveyResultsRepo repository.SurveyResults
}

func NewSurveysService(modulesRepo repository.Modules, surveyResultsRepo repository.SurveyResults) *SurveysService {
	return &SurveysService{modulesRepo: modulesRepo, surveyResultsRepo: surveyResultsRepo}
}

func (s *SurveysService) Create(ctx context.Context, inp CreateSurveyInput) error {
	for i := range inp.Survey.Questions {
		inp.Survey.Questions[i].ID = primitive.NewObjectID()
	}

	return s.modulesRepo.AttachSurvey(ctx, inp.SchoolID, inp.ModuleID, inp.Survey)
}

func (s *SurveysService) Delete(ctx context.Context, schoolId, moduleId primitive.ObjectID) error {
	return s.modulesRepo.DetachSurvey(ctx, schoolId, moduleId)
}

func (s *SurveysService) SaveStudentAnswers(ctx context.Context, inp SaveStudentAnswersInput) error {
	return s.surveyResultsRepo.Save(ctx, domain.SurveyResult{
		StudentID:   inp.StudentID,
		ModuleID:    inp.ModuleID,
		SubmittedAt: time.Now(),
		Answers:     inp.Answers,
	})
}

func (s *SurveysService) GetResultsByModule(ctx context.Context, moduleId primitive.ObjectID,
	pagination *domain.PaginationQuery) ([]domain.SurveyResult, int64, error) {
	return s.surveyResultsRepo.GetAllByModule(ctx, moduleId, pagination)
}
