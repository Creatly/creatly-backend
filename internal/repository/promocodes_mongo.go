package repository

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
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

func (r *PromocodesRepo) GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.PromoCode, error) {
	var promocode domain.PromoCode
	if err := r.db.FindOne(ctx, bson.M{"schoolId": schoolId, "code": code}).Decode(&promocode); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.PromoCode{}, ErrPromoNotFound
		}

		return domain.PromoCode{}, err
	}

	return promocode, nil
}

func (r *PromocodesRepo) GetById(ctx context.Context, schoolId, id primitive.ObjectID) (domain.PromoCode, error) {
	var promocode domain.PromoCode
	if err := r.db.FindOne(ctx, bson.M{"_id": id, "schoolId": schoolId}).Decode(&promocode); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.PromoCode{}, ErrPromoNotFound
		}

		return domain.PromoCode{}, err
	}

	return promocode, nil
}
