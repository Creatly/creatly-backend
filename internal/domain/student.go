package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Student struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	RegisteredAt       int64                `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt        int64                `json:"lastVisitAt" bson:"lastVisitAt"`
	SchoolID           primitive.ObjectID   `json:"schoolId" bson:"schoolId"`
	SourceCourseID     primitive.ObjectID   `json:"sourceCourseId" bson:"sourceCourseId"`
	Orders             []Order              `json:"orders" bson:"orders"`
	AvailableModuleIDs []primitive.ObjectID `json:"availableModuleIds" bson:"availableModuleIds"`
	Verification       Verification         `json:"verification" bson:"verification"`
}

type Verification struct {
	Hash     string `json:"hash" bson:"hash"`
	Verified bool   `json:"verified" bson:"verified"`
}
