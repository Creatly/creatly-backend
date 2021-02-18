package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Offer struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name"`
	Description string               `json:"description" bson:"description,omitempty"`
	SchoolID    primitive.ObjectID   `json:"schoolId" bson:"schoolId"`
	PackageIDs  []primitive.ObjectID `json:"packages" bson:"packages,omitempty"`
	Price       Price                `json:"price" bson:"price"`
}

type Price struct {
	Value    uint    `json:"value" bson:"value"`
	Currency string `json:"currency" bson:"currency"`
}
