package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type School struct {
	ID           primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name         string              `json:"name" bson:"name"`
	Description  string              `json:"description" bson:"description"`
	RegisteredAt primitive.Timestamp `json:"registeredAt" bson:"registeredAt"`
	Admins       []Admin             `json:"admins" bson:"admins"`
	Courses      []Course            `json:"courses" bson:"courses"`
	Settings     Settings            `json:"settings" bson:"settings"`
}

type Settings struct {
	Color       string `json:"color" bson:"color"`
	Domain      string `json:"domain" bson:"domain"`
	Email       string `json:"email" bson:"email"`
	ContactData string `json:"contactData" bson:"contactData"`
	Pages       Pages  `json:"pages" bson:"pages"`
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
