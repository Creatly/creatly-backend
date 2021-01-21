package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Offer struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description"`
	CreatedAt   int64                `json:"createdAt" bson:"createdAt"`
	SchoolID    primitive.ObjectID   `json:"schoolId" bson:"schoolId"`
	PackageIDs  []primitive.ObjectID `json:"packageIds" bson:"packageIds"`
	Price       Price                `json:"price" bson:"price"`
}

type Price struct {
	Value    float64 `json:"value" bson:"value"` // TODO store in int?
	Currency string  `json:"currency" bson:"currency"`
}
