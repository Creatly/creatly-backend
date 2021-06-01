package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type School struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description,omitempty"`
	RegisteredAt time.Time          `json:"registeredAt" bson:"registeredAtm,omitempty"`
	Admins       []Admin            `json:"admins" bson:"admins,omitempty"`
	Courses      []Course           `json:"courses" bson:"courses,omitempty"`
	Settings     Settings           `json:"settings" bson:"settings,omitempty"`
}

type Settings struct {
	Color       string      `json:"color" bson:"color,omitempty"`
	Domains     []string    `json:"domains" bson:"domains,omitempty"`
	ContactInfo ContactInfo `json:"contactInfo" bson:"contactInfo,omitempty"`
	Pages       Pages       `json:"pages" bson:"pages,omitempty"`
}

// todo review fields.
type ContactInfo struct {
	BusinessName       string `json:"businessName" bson:"businessName,omitempty"`
	RegistrationNumber string `json:"registrationNumber" bson:"registrationNumber,omitempty"`
	Address            string `json:"address" bson:"address,omitempty"`
	Email              string `json:"email" bson:"email,omitempty"`
}

type Pages struct {
	Confidential     string `json:"confidential" bson:"confidential,omitempty"`
	ServiceAgreement string `json:"serviceAgreement" bson:"serviceAgreement,omitempty"`
	RefundPolicy     string `json:"refundPolicy" bson:"refundPolicy,omitempty"`
}

type Admin struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	SchoolId primitive.ObjectID
}
