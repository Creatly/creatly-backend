package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	OrderStatusCreated  = "created"
	OrderStatusPaid     = "paid"
	OrderStatusFailed   = "failed"
	OrderStatusCanceled = "canceled"
	OrderStatusOther    = "other"
)

type Order struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SchoolID     primitive.ObjectID `json:"schoolId" bson:"schoolId"`
	Student      StudentInfoShort   `json:"student" bson:"student"`
	Offer        OrderOfferInfo     `json:"offer" bson:"offer"`
	Promo        OrderPromoInfo     `json:"promo" bson:"promo,omitempty"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	Amount       uint               `json:"amount" bson:"amount"`
	Currency     string             `json:"currency" bson:"currency"`
	Status       string             `json:"status" bson:"status"`
	Transactions []Transaction      `json:"transactions" bson:"transactions,omitempty"`
}

type OrderOfferInfo struct {
	ID   primitive.ObjectID `json:"id" bson:"id"`
	Name string             `json:"name" bson:"name"`
}

type OrderPromoInfo struct {
	ID   primitive.ObjectID `json:"id" bson:"id"`
	Code string             `json:"code" bson:"code"`
}

type Transaction struct {
	Status         string    `json:"status" bson:"status"`
	CreatedAt      time.Time `json:"createdAt" bson:"createdAt"`
	AdditionalInfo string    `json:"additionalInfo" bson:"additionalInfo"`
}
