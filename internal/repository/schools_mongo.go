package repository

import (
	"context"
	"time"

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

func (r *SchoolsRepo) Create(ctx context.Context, name string) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, domain.School{
		Name:         name,
		RegisteredAt: time.Now(),
	})

	return res.InsertedID.(primitive.ObjectID), err
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

	if inp.Color != nil {
		updateQuery["settings.color"] = inp.Color
	}

	if inp.Domains != nil {
		updateQuery["settings.domains"] = inp.Domains
	}

	if inp.Email != nil {
		updateQuery["settings.email"] = inp.Email
	}

	if inp.ContactInfo != nil {
		updateQuery["settings.contactInfo"] = inp.ContactInfo
	}

	if inp.Pages != nil {
		updateQuery["settings.pages"] = inp.Pages
	}

	if inp.ShowPaymentImages != nil {
		updateQuery["settings.showPaymentImages"] = inp.ShowPaymentImages
	}

	if inp.GoogleAnalyticsCode != nil {
		updateQuery["settings.googleAnalyticsCode"] = *inp.GoogleAnalyticsCode
	}

	if inp.LogoURL != nil {
		updateQuery["settings.logo"] = *inp.LogoURL
	}

	_, err := r.db.UpdateOne(ctx,
		bson.M{"_id": inp.SchoolID}, bson.M{"$set": updateQuery})

	return err
}

func (r *SchoolsRepo) SetFondyCredentials(ctx context.Context, id primitive.ObjectID, fondy domain.Fondy) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"settings.fondy": fondy}})

	return err
}
