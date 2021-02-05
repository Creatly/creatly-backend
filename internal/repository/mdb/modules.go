package mdb

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ModulesRepo struct {
	db *mongo.Collection
}

func NewModulesRepo(db *mongo.Database) *ModulesRepo {
	return &ModulesRepo{db: db.Collection(modulesCollection)}
}

func (r *ModulesRepo) Create(ctx context.Context, module domain.Module) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, module)
	return res.InsertedID.(primitive.ObjectID), err
}

func (r *ModulesRepo) GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Module, error) {
	var modules []domain.Module
	cur, err := r.db.Find(ctx, bson.M{"courseId": courseId})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &modules)
	return modules, err
}

func (r *ModulesRepo) GetById(ctx context.Context, moduleId primitive.ObjectID) (domain.Module, error) {
	var module domain.Module
	err := r.db.FindOne(ctx, bson.M{"_id": moduleId}).Decode(&module)
	return module, err
}

func (r *ModulesRepo) GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Module, error) {
	var modules []domain.Module
	cur, err := r.db.Find(ctx, bson.M{"packageId": bson.M{"$in": packageIds}})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &modules)
	return modules, err
}
