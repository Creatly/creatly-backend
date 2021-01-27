package mdb

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrdersRepo struct {
	db *mongo.Collection
}

func NewOrdersRepo(db *mongo.Database) *OrdersRepo {
	return &OrdersRepo{
		db: db.Collection(ordersCollection),
	}
}

func (r *OrdersRepo) Create(ctx context.Context, order domain.Order) error {
	_, err := r.db.InsertOne(ctx, order)
	return err
}

func (r *OrdersRepo) AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) error {
	res := r.db.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"status": transaction.Status,
		},
		"$push": bson.M{
			"transactions": transaction,
		},
	})
	return res.Err()
}
