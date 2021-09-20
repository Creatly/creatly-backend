package repository

import (
	"context"

	"github.com/zhashkevych/creatly-backend/internal/domain"
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
	var packages []domain.Package

	cur, err := r.db.Find(ctx, bson.M{"courseId": courseId})
	if err != nil {
		return packages, err
	}

	err = cur.All(ctx, &packages)

	return packages, err
}

func (r *PackagesRepo) GetById(ctx context.Context, id primitive.ObjectID) (domain.Package, error) {
	var pkg domain.Package
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&pkg)

	return pkg, err
}

func (r *PackagesRepo) GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.Package, error) {
	var pkgs []domain.Package

	cur, err := r.db.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &pkgs)

	return pkgs, err
}

func (r *PackagesRepo) Update(ctx context.Context, inp UpdatePackageInput) error {
	updateQuery := bson.M{}

	if inp.Name != "" {
		updateQuery["name"] = inp.Name
	}

	_, err := r.db.UpdateOne(ctx,
		bson.M{"_id": inp.ID, "schoolId": inp.SchoolID}, bson.M{"$set": updateQuery})

	return err
}

func (r *PackagesRepo) Delete(ctx context.Context, schoolId, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id, "schoolId": schoolId})

	return err
}
