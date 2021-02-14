package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID               primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name             string               `json:"name" bson:"name"`
	Email            string               `json:"email" bson:"email"`
	Password         string               `json:"password" bson:"password"`
	RegisteredAt     time.Time            `json:"registeredAt" bson:"registeredAt"`
	LastVisitAt      time.Time            `json:"lastVisitAt" bson:"lastVisitAt"`
	SchoolID         primitive.ObjectID   `json:"schoolId" bson:"schoolId"`
	AvailableModules []primitive.ObjectID `json:"availableModules" bson:"availableModules,omitempty"`
	AvailableCourses []primitive.ObjectID `json:"availableCourses" bson:"availableCourses,omitempty"`
	Verification     Verification         `json:"verification" bson:"verification"`
	Session          Session              `json:"session" bson:"session,omitempty"`
}

func (s Student) IsModuleAvailable(m Module) bool {
	for _, id := range s.AvailableModules {
		if m.ID == id {
			return true
		}
	}
	return false
}

type Verification struct {
	Code     string `json:"code" bson:"code"`
	Verified bool   `json:"verified" bson:"verified"`
}
