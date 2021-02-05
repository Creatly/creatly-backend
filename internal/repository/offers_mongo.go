package repository

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
	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&offer)
	return offer, err
}
