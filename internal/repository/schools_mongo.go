package repository

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SchoolsRepo struct {
	db *mongo.Collection
}

func NewSchoolsRepo(db *mongo.Database) *SchoolsRepo {
	return &SchoolsRepo{
		db: db.Collection(schoolsCollection),
	}
}

func (r *SchoolsRepo) GetByDomain(ctx context.Context, domainName string) (domain.School, error) {
	var school domain.School
	err := r.db.FindOne(ctx, bson.M{"domain": domainName}).Decode(&school)

	return school, err
}

func (r *SchoolsRepo) GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error) {
	var school domain.School
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&school)

	return school, err
}
