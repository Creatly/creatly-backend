package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type BaseUser struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	RegisteredAt int64              `json:"registeredAt" bson:"registeredAt"`
	Schools      []School           `json:"schools" bson:"schools"`
	BaseUser
}
