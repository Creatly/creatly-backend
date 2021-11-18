package repository

import (
	"context"
	"errors"

	"github.com/zhashkevych/creatly-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PromocodesRepo struct {
	db *mongo.Collection
}

func NewPromocodeRepo(db *mongo.Database) *PromocodesRepo {
	return &PromocodesRepo{db: db.Collection(promocodesCollection)}
}

func (r *PromocodesRepo) Create(ctx context.Context, promocode domain.PromoCode) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, promocode)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *PromocodesRepo) Update(ctx context.Context, inp domain.UpdatePromoCodeInput) error {
	updateQuery := bson.M{}

	if inp.Code != "" {
		updateQuery["code"] = inp.Code
	}

	if inp.DiscountPercentage != 0 {
		updateQuery["discountPercentage"] = inp.DiscountPercentage
	}

	if !inp.ExpiresAt.IsZero() {
		updateQuery["expiresAt"] = inp.ExpiresAt
	}

	if inp.OfferIDs != nil {
		updateQuery["offerIds"] = inp.OfferIDs
	}

	_, err := r.db.UpdateOne(ctx,
		bson.M{"_id": inp.ID, "schoolId": inp.SchoolID}, bson.M{"$set": updateQuery})

	return err
}

func (r *PromocodesRepo) Delete(ctx context.Context, schoolId, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": id, "schoolId": schoolId})

	return err
}

func (r *PromocodesRepo) GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error) {
	var promocode domain.PromoCode
	if err := r.db.FindOne(ctx, bson.M{"schoolId": schoolId, "code": code}).Decode(&promocode); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.PromoCode{}, domain.ErrPromoNotFound
		}

		return domain.PromoCode{}, err
	}

	return promocode, nil
}

func (r *PromocodesRepo) GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error) {
	var promocode domain.PromoCode
	if err := r.db.FindOne(ctx, bson.M{"_id": id, "schoolId": schoolId}).Decode(&promocode); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.PromoCode{}, domain.ErrPromoNotFound
		}

		return domain.PromoCode{}, err
	}

	return promocode, nil
}

func (r *PromocodesRepo) GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.PromoCode, error) {
	cursor, err := r.db.Find(ctx, bson.M{"schoolId": schoolId})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrPromoNotFound
		}

		return nil, err
	}

	var promocodes []domain.PromoCode
	if err = cursor.All(ctx, &promocodes); err != nil {
		return nil, err
	}

	return promocodes, nil
}
