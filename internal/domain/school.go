package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type School struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Description  string             `json:"description" bson:"description"`
	RegisteredAt time.Time          `json:"registeredAt" bson:"registeredAt"`
	Admins       []Admin            `json:"admins" bson:"admins"`
	Courses      []Course           `json:"courses" bson:"courses"`
	Settings     Settings           `json:"settings" bson:"settings"`
}

type Settings struct {
	Color       string      `json:"color" bson:"color"`
	Domains     []string    `json:"domains" bson:"domains"`
	ContactInfo ContactInfo `json:"contactInfo" bson:"contactInfo"`
	Pages       Pages       `json:"pages" bson:"pages"`
}

type ContactInfo struct {
	BusinessName       string `json:"businessName" bson:"businessName"`
	RegistrationNumber string `json:"registrationNumber" bson:"RegistrationNumber"`
	Address            string `json:"address" bson:"address"`
	Email              string `json:"email" bson:"email"`
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
