package repository

import (
	"context"
	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *SurveyResultsRepo) GetAllByModule(ctx context.Context, moduleId primitive.ObjectID, pagination *domain.PaginationQuery) ([]domain.SurveyResult, int64, error) {
	opts := getPaginationOpts(pagination)
	filter := bson.M{"moduleId": moduleId}

	cur, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var results []domain.SurveyResult
	if err := cur.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	count, err := r.db.CountDocuments(ctx, filter)

	return results, count, err
}
