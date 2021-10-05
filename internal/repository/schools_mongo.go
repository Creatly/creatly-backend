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

func (r *SchoolsRepo) UpdateSettings(ctx context.Context, id primitive.ObjectID, inp domain.UpdateSchoolSettingsInput) error {
	updateQuery := bson.M{}

	if inp.Name != nil {
		updateQuery["name"] = *inp.Name
	}

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
		setContactInfoUpdateQuery(&updateQuery, inp)
	}

	if inp.Pages != nil {
		setPagesUpdateQuery(&updateQuery, inp)
	}

	if inp.ShowPaymentImages != nil {
		updateQuery["settings.showPaymentImages"] = inp.ShowPaymentImages
	}

	if inp.DisableRegistration != nil {
		updateQuery["settings.disableRegistration"] = inp.DisableRegistration
	}

	if inp.GoogleAnalyticsCode != nil {
		updateQuery["settings.googleAnalyticsCode"] = *inp.GoogleAnalyticsCode
	}

	if inp.LogoURL != nil {
		updateQuery["settings.logo"] = *inp.LogoURL
	}

	_, err := r.db.UpdateOne(ctx,
		bson.M{"_id": id}, bson.M{"$set": updateQuery})

	return err
}

func (r *SchoolsRepo) SetFondyCredentials(ctx context.Context, id primitive.ObjectID, fondy domain.Fondy) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"settings.fondy": fondy}})

	return err
}

func setContactInfoUpdateQuery(updateQuery *bson.M, inp domain.UpdateSchoolSettingsInput) {
	if inp.ContactInfo.Address != nil {
		(*updateQuery)["settings.contactInfo.address"] = inp.ContactInfo.Address
	}

	if inp.ContactInfo.BusinessName != nil {
		(*updateQuery)["settings.contactInfo.businessName"] = inp.ContactInfo.BusinessName
	}

	if inp.ContactInfo.Email != nil {
		(*updateQuery)["settings.contactInfo.email"] = inp.ContactInfo.Email
	}

	if inp.ContactInfo.Phone != nil {
		(*updateQuery)["settings.contactInfo.phone"] = inp.ContactInfo.Phone
	}

	if inp.ContactInfo.RegistrationNumber != nil {
		(*updateQuery)["settings.contactInfo.registrationNumber"] = inp.ContactInfo.RegistrationNumber
	}
}

func setPagesUpdateQuery(updateQuery *bson.M, inp domain.UpdateSchoolSettingsInput) {
	if inp.Pages.Confidential != nil {
		(*updateQuery)["settings.pages.confidential"] = inp.Pages.Confidential
	}

	if inp.Pages.NewsletterConsent != nil {
		(*updateQuery)["settings.pages.newsletterConsent"] = inp.Pages.NewsletterConsent
	}

	if inp.Pages.ServiceAgreement != nil {
		(*updateQuery)["settings.pages.serviceAgreement"] = inp.Pages.ServiceAgreement
	}
}
