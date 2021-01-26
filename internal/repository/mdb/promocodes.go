package mdb

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

func (r *PromocodesRepo) GetByCode(ctx context.Context, schoolId primitive.ObjectID, code string) (domain.Promocode, error) {
	var promocode domain.Promocode
	err := r.db.FindOne(ctx, bson.M{"schoolId": schoolId, "code": code}).Decode(&promocode)
	return promocode, err
}

func (r *PromocodesRepo) GetById(ctx context.Context, id primitive.ObjectID) (domain.Promocode, error) {
	var promocode domain.Promocode
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&promocode)
	return promocode, err
}
