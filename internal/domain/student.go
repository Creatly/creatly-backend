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
	AvailableOffers  []primitive.ObjectID `json:"availableOffers" bson:"availableOffers,omitempty"`
	Verification     Verification         `json:"verification" bson:"verification"`
	Session          Session              `json:"session" bson:"session,omitempty"`
	Blocked          bool                 `json:"blocked" bson:"blocked"`
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

type StudentLessons struct {
	StudentID  primitive.ObjectID   `json:"studentId" bson:"studentId"`
	Finished   []primitive.ObjectID `json:"finished" bson:"finished"`
	LastOpened primitive.ObjectID   `json:"lastOpened" bson:"lastOpened"`
}

type StudentInfoShort struct {
	ID    primitive.ObjectID `json:"id" bson:"id"`
	Name  string             `json:"name" bson:"name"`
	Email string             `json:"email" bson:"email"`
}

type UpdateStudentInput struct {
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Verified  *bool              `json:"verified"`
	Blocked   *bool              `json:"blocked"`
	StudentID primitive.ObjectID `json:"-"`
	SchoolID  primitive.ObjectID `json:"-"`
}

type CreateStudentInput struct {
	Name     string             `json:"name" binding:"required,min=2"`
	Email    string             `json:"email" binding:"required,email"`
	Password string             `json:"password" binding:"required,min=6"`
	SchoolID primitive.ObjectID `json:"-"`
}
