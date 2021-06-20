package repository

import (
	"context"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	if err := cur.All(ctx, &content); err != nil {
		return nil, err
	}

	return content, nil
}

func (r *LessonContentRepo) GetByLesson(ctx context.Context, lessonId primitive.ObjectID) (domain.LessonContent, error) {
	var content domain.LessonContent
	err := r.db.FindOne(ctx, bson.M{"lessonId": lessonId}).Decode(&content)

	return content, err
}

func (r *LessonContentRepo) Update(ctx context.Context, schoolId, lessonId primitive.ObjectID, content string) error {
	opts := &options.UpdateOptions{}
	opts.SetUpsert(true)

	_, err := r.db.UpdateOne(ctx, bson.M{"lessonId": lessonId, "schoolId": schoolId}, bson.M{"$set": bson.M{"content": content}}, opts)

	return err
}

func (r *LessonContentRepo) DeleteContent(ctx context.Context, schoolId primitive.ObjectID, lessonIds []primitive.ObjectID) error {
	_, err := r.db.DeleteMany(ctx, bson.M{"lessonId": bson.M{"$in": lessonIds}, "schoolId": schoolId})

	return err
}
