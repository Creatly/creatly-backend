package service

import (
	"context"
	"time"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"github.com/zhashkevych/creatly-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SurveysService struct {
	modulesRepo       repository.Modules
	surveyResultsRepo repository.SurveyResults
	studentsRepo      repository.Students
}

func NewSurveysService(modulesRepo repository.Modules, surveyResultsRepo repository.SurveyResults, studentsRepo repository.Students) *SurveysService {
	return &SurveysService{modulesRepo: modulesRepo, surveyResultsRepo: surveyResultsRepo, studentsRepo: studentsRepo}
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
	student, err := s.studentsRepo.GetById(ctx, inp.SchoolID, inp.StudentID)
	if err != nil {
		return err
	}

	return s.surveyResultsRepo.Save(ctx, domain.SurveyResult{
		Student: domain.StudentInfoShort{
			ID:    student.ID,
			Name:  student.Name,
			Email: student.Email,
		},
		ModuleID:    inp.ModuleID,
		SubmittedAt: time.Now(),
		Answers:     inp.Answers,
	})
}

func (s *SurveysService) GetResultsByModule(ctx context.Context, moduleId primitive.ObjectID,
	pagination *domain.PaginationQuery) ([]domain.SurveyResult, int64, error) {
	return s.surveyResultsRepo.GetAllByModule(ctx, moduleId, pagination)
}

func (s *SurveysService) GetStudentResults(ctx context.Context, moduleID, studentID primitive.ObjectID) (domain.SurveyResult, error) {
	return s.surveyResultsRepo.GetByStudent(ctx, moduleID, studentID)
}
