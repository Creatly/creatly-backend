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

func (r *OrdersRepo) GetBySchool(ctx context.Context, schoolId primitive.ObjectID, query domain.GetOrdersQuery) ([]domain.Order, int64, error) {
	opts := getPaginationOpts(&query.PaginationQuery)
	opts.SetSort(bson.M{"createdAt": -1})

	filter := bson.M{"$and": []bson.M{{"schoolId": schoolId}}}

	if query.Search != "" {
		expression := primitive.Regex{Pattern: query.Search}

		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			"$or": []bson.M{
				{"student.name": expression},
				{"student.email": expression},
				{"offer.name": expression},
				{"promo.name": expression},
			},
		})
	}

	if query.Status != "" {
		filter["$and"] = append(filter["$and"].([]bson.M), bson.M{
			"status": query.Status,
		})
	}

	if err := filterDateQueries(query.DateFrom, query.DateTo, "createdAt", filter); err != nil {
		return nil, 0, err
	}

	cur, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var orders []domain.Order
	if err := cur.All(ctx, &orders); err != nil {
		return nil, 0, err
	}

	count, err := r.db.CountDocuments(ctx, filter)

	return orders, count, err
}

func (r *OrdersRepo) GetById(ctx context.Context, id primitive.ObjectID) (domain.Order, error) {
	var order domain.Order

	err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&order)

	return order, err
}

func (r *OrdersRepo) SetStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})

	return err
}
