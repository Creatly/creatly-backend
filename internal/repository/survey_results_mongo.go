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

func (r *SurveyResultsRepo) GetAllByModule(ctx context.Context, moduleID primitive.ObjectID, pagination *domain.PaginationQuery) ([]domain.SurveyResult, int64, error) {
	opts := getPaginationOpts(pagination)
	filter := bson.M{"moduleId": moduleID}

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

func (r *SurveyResultsRepo) GetByStudent(ctx context.Context, moduleID, studentID primitive.ObjectID) (domain.SurveyResult, error) {
	var res domain.SurveyResult
	err := r.db.FindOne(ctx, bson.M{"student.id": studentID, "moduleId": moduleID}).Decode(&res)

	return res, err
}
