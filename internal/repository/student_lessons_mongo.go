package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StudentLessonsRepo struct {
	db *mongo.Collection
}

func NewStudentLessonsRepo(db *mongo.Database) *StudentLessonsRepo {
	return &StudentLessonsRepo{db: db.Collection(studentLessonsCollection)}
}

func (r *StudentLessonsRepo) AddFinished(ctx context.Context, studentId, lessonId primitive.ObjectID) error {
	filter := bson.M{"studentId": studentId}
	update := bson.M{"$addToSet": bson.M{"finished": lessonId}}

	_, err := r.db.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

	return err
}

func (r *StudentLessonsRepo) SetLastOpened(ctx context.Context, studentId, lessonId primitive.ObjectID) error {
	filter := bson.M{"studentId": studentId}
	update := bson.M{"$set": bson.M{"lastOpened": lessonId}}

	_, err := r.db.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

	return err
}
