package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Offer struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description,omitempty"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
	SchoolID    primitive.ObjectID   `json:"schoolId" bson:"schoolId"`
	PackageIDs  []primitive.ObjectID `json:"packageIds" bson:"packageIds"`
	Price       Price                `json:"price" bson:"price"`
}

type Price struct {
	Value    int    `json:"value" bson:"value"`
	Currency string `json:"currency" bson:"currency"`
}
