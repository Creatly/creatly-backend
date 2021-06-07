package repository

import (
	"context"

	"github.com/zhashkevych/creatly-backend/internal/domain"
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

func (r *OrdersRepo) AddTransaction(ctx context.Context, id primitive.ObjectID, transaction domain.Transaction) (domain.Order, error) {
	var order domain.Order

	res := r.db.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"status": transaction.Status,
		},
		"$push": bson.M{
			"transactions": transaction,
		},
	})
	if res.Err() != nil {
		return order, res.Err()
	}

	err := res.Decode(&order)

	return order, err
}

func (r *OrdersRepo) GetBySchool(ctx context.Context, schoolId primitive.ObjectID) ([]domain.Order, error) {
	cur, err := r.db.Find(ctx, bson.M{"schoolId": schoolId})
	if err != nil {
		return nil, err
	}

	var orders []domain.Order
	err = cur.All(ctx, &orders)

	return orders, err
}
