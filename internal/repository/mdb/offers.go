package mdb

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
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

func (r *OffersRepo) GetSchoolOffers(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Offer, error) {
	cur, err := r.db.Find(ctx, bson.M{"schoolId": schoolId})
	if err != nil {
		return nil, err
	}

	var offers []domain.Offer
	err = cur.All(ctx, &offers)
	return offers, err
}
