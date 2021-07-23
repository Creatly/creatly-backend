package repository

import (
	"context"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type SurveyResultsRepo struct {
	db *mongo.Collection
}

func NewSurveyResultsRepo(db *mongo.Database) *SurveyResultsRepo {
	return &SurveyResultsRepo{
		db: db.Collection(surveyResultsCollection),
	}
}

func (r *SurveyResultsRepo) Save(ctx context.Context, results domain.SurveyResult) error {
	_, err := r.db.InsertOne(ctx, results)

	return err
}
