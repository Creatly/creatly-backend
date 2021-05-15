package repository

import (
	"context"

	"github.com/zhashkevych/creatly-backend/internal/domain"
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
	err := r.db.FindOne(ctx, bson.M{"settings.domains": domainName}).Decode(&school)

	return school, err
}

func (r *SchoolsRepo) GetById(ctx context.Context, id primitive.ObjectID) (domain.School, error) {
	var school domain.School
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&school)

	return school, err
}

func (r *SchoolsRepo) UpdateSettings(ctx context.Context, inp UpdateSchoolSettingsInput) error {
	updateQuery := bson.M{}

	if inp.Color != "" {
		updateQuery["settings.color"] = inp.Color
	}

	if inp.Domains != nil {
		updateQuery["settings.domains"] = inp.Domains
	}

	if inp.Email != "" {
		updateQuery["settings.email"] = inp.Email
	}

	if inp.ContactInfo != nil {
		updateQuery["settings.contactInfo"] = inp.ContactInfo
	}

	if inp.Pages != nil {
		updateQuery["settings.pages"] = inp.Pages
	}

	_, err := r.db.UpdateOne(ctx,
		bson.M{"_id": inp.SchoolID}, bson.M{"$set": updateQuery})
	return err
}
