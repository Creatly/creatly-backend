package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrFondyIsNotConnected = errors.New("fondy is not connected")

type School struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Subtitle     string             `json:"subtitle" bson:"subtitle,omitempty"`
	Description  string             `json:"description" bson:"description,omitempty"`
	RegisteredAt time.Time          `json:"registeredAt" bson:"registeredAtm,omitempty"`
	Admins       []Admin            `json:"admins" bson:"admins,omitempty"`
	Courses      []Course           `json:"courses" bson:"courses,omitempty"`
	Settings     Settings           `json:"settings" bson:"settings,omitempty"`
}

type Settings struct {
	Color               string      `json:"color" bson:"color,omitempty"`
	Domains             []string    `json:"domains" bson:"domains,omitempty"`
	ContactInfo         ContactInfo `json:"contactInfo" bson:"contactInfo,omitempty"`
	Pages               Pages       `json:"pages" bson:"pages,omitempty"`
	ShowPaymentImages   bool        `json:"showPaymentImages" bson:"showPaymentImages,omitempty"`
	Logo                string      `json:"logo" bson:"logo,omitempty"`
	GoogleAnalyticsCode string      `json:"googleAnalyticsCode" bson:"googleAnalyticsCode,omitempty"`
	Fondy               Fondy       `json:"fondy" bson:"fondy,omitempty"`
}

type Fondy struct {
	MerchantID       string `json:"merchantId" bson:"merchantId"`
	MerchantPassword string `json:"merchantPassword" bson:"merchantPassword"`
	Connected        bool   `json:"connected" bson:"connected"`
}

type ContactInfo struct {
	BusinessName       string `json:"businessName" bson:"businessName,omitempty"`
	RegistrationNumber string `json:"registrationNumber" bson:"registrationNumber,omitempty"`
	Address            string `json:"address" bson:"address,omitempty"`
	Email              string `json:"email" bson:"email,omitempty"`
	Phone              string `json:"phone" bson:"phone,omitempty"`
}

type Pages struct {
	Confidential      string `json:"confidential" bson:"confidential,omitempty"`
	ServiceAgreement  string `json:"serviceAgreement" bson:"serviceAgreement,omitempty"`
	RefundPolicy      string `json:"refundPolicy" bson:"refundPolicy,omitempty"`
	NewsletterConsent string `json:"newsletterConsent" bson:"newsletterConsent,omitempty"`
}

type Admin struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	SchoolID primitive.ObjectID
}
