package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"hash" bson:"hash"`
	RegisteredAt int64              `json:"registeredAt" bson:"registeredAt"`
	Schools      []School           `json:"schools" bson:"schools"`
}
