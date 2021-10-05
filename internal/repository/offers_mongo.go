package repository

import (
	"context"
	"errors"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OffersRepo struct {
	db *mongo.Collection
}

func NewOffersRepo(db *mongo.Database) *OffersRepo {
	return &OffersRepo{
		db: db.Collection(offersCollection),
	}
}

func (r *OffersRepo) GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Offer, error) {
	cur, err := r.db.Find(ctx, bson.M{"schoolId": schoolId})
	if err != nil {
		return nil, err
	}

	var offers []domain.Offer
	err = cur.All(ctx, &offers)

	return offers, err
}

func (r *OffersRepo) GetById(ctx context.Context, id primitive.ObjectID) (domain.Offer, error) {
	var offer domain.Offer
	if err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&offer); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Offer{}, domain.ErrOfferNotFound
		}

		return domain.Offer{}, err
	}

	return offer, nil
}

func (r *OffersRepo) GetByPackages(ctx context.Context, packageIds []primitive.ObjectID) ([]domain.Offer, error) {
	cur, err := r.db.Find(ctx, bson.M{"packages": bson.M{"$in": packageIds}})
	if err != nil {
		return nil, err
	}

	var offers []domain.Offer
	err = cur.All(ctx, &offers)

	return offers, err
}

func (r *OffersRepo) Create(ctx context.Context, offer domain.Offer) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, offer)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *OffersRepo) Update(ctx context.Context, inp UpdateOfferInput) error {
	updateQuery := bson.M{}

	if inp.Name != "" {
		updateQuery["name"] = inp.Name
	}

	if inp.Description != "" {
		updateQuery["description"] = inp.Description
	}

	if inp.Benefits != nil {
		updateQuery["benefits"] = inp.Benefits
	}

	if inp.Price != nil {
		updateQuery["price"] = inp.Price
	}

	if inp.Packages != nil {
		updateQuery["packages"] = inp.Packages
	}

	if inp.PaymentMethod != nil {
		updateQuery["paymentMethod"] = inp.PaymentMethod
	}

	_, err := r.db.UpdateOne(ctx,
		bson.M{"_id": inp.ID, "schoolId": inp.SchoolID}, bson.M{"$set": updateQuery})

	return err
}

func (r *OffersRepo) Delete(ctx context.Context, schoolId, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id, "schoolId": schoolId})

	return err
}

func (r OffersRepo) GetByIds(ctx context.Context, ids []primitive.ObjectID) ([]domain.Offer, error) {
	var offers []domain.Offer

	cur, err := r.db.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &offers)

	return offers, err
}
