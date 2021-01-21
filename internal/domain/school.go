package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type School struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	Domain       string             `json:"domain" bson:"domain"`
	RegisteredAt int64              `json:"registeredAt" bson:"registeredAt"`
	Admins       []Admin            `json:"admins" bson:"admins"`
	Courses      []Course           `json:"courses" bson:"courses"`
}

type Admin struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	BaseUser
}
