package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Student struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name               string               `json:"name" bson:"name"`
	Email              string               `json:"email" bson:"email"`
	Password           string               `json:"password" bson:"password"`
	RegisteredAt       time.Time            `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt        time.Time            `json:"lastVisitAt" bson:"lastVisitAt"`
	SchoolID           primitive.ObjectID   `json:"schoolId" bson:"schoolId"`
	SourceCourseID     primitive.ObjectID   `json:"sourceCourseId" bson:"sourceCourseId,omitempty"`
	Orders             []Order              `json:"orders" bson:"orders,omitempty"`
	AvailableModuleIDs []primitive.ObjectID `json:"availableModuleIds" bson:"availableModuleIds,omitempty"`
	Verification       Verification         `json:"verification" bson:"verification"`
}

type Verification struct {
	Hash     primitive.ObjectID `json:"hash" bson:"hash"`
	Verified bool               `json:"verified" bson:"verified"`
}
