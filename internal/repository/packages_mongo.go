package repository

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PackagesRepo struct {
	db *mongo.Collection
}

func NewPackagesRepo(db *mongo.Database) *PackagesRepo {
	return &PackagesRepo{db: db.Collection(packagesCollection)}
}

func (r *PackagesRepo) Create(ctx context.Context, pkg domain.Package) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, pkg)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *PackagesRepo) GetByCourse(ctx context.Context, courseId primitive.ObjectID) ([]domain.Package, error) {
	var pkgs []domain.Package
	cur, err := r.db.Find(ctx, bson.M{"courseId": courseId})
	if err != nil {
		return pkgs, err
	}

	err = cur.All(ctx, &pkgs)
	return pkgs, err
}
