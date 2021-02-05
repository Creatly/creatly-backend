package mdb

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LessonContentRepo struct {
	db *mongo.Collection
}

func NewLessonContentRepo(db *mongo.Database) *LessonContentRepo {
	return &LessonContentRepo{db: db.Collection(contentCollection)}
}

func (r *LessonContentRepo) GetByLessons(ctx context.Context, lessonIds []primitive.ObjectID) ([]domain.LessonContent, error) {
	var content []domain.LessonContent
	cur, err := r.db.Find(ctx, bson.M{"lessonId": bson.M{"$in": lessonIds}})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}
