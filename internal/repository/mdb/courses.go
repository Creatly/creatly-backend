package mdb

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CoursesRepo struct {
	db *mongo.Database
}

func NewCoursesRepo(db *mongo.Database) *CoursesRepo {
	return &CoursesRepo{db: db}
}

func (r *CoursesRepo) GetModules(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	var modules []domain.Module
	cur, err := r.db.Collection(modulesCollection).Find(ctx, bson.M{"courseId": courseId, "published": true})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &modules)
	return modules, err
}
