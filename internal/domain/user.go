package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name         string               `json:"name" bson:"name"`
	Email        string               `json:"email" bson:"email"`
	Phone        string               `json:"phone" bson:"phone"`
	Password     string               `json:"password" bson:"password"`
	RegisteredAt time.Time            `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt  time.Time            `json:"lastVisitAt" bson:"lastVisitAt"`
	Verification Verification         `json:"verification" bson:"verification"`
	Schools      []primitive.ObjectID `json:"schools" bson:"schools"`
}
