package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name         string               `json:"name" bson:"name"`
	Email        string               `json:"email" bson:"email"`
	Password     string               `json:"password" bson:"password"`
	RegisteredAt int64                `json:"registeredAt" bson:"registeredAt"`
	Schools      []primitive.ObjectID `json:"schools" bson:"schools"`
}
