package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	OrderStatusCreated = "created"
	OrderStatusPaid    = "paid"
	OrderStatusFailed  = "failed"
	OrderStatusOther   = "other"
)

type Order struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StudentID    primitive.ObjectID `json:"studentId" bson:"studentId"`
	OfferID      primitive.ObjectID `json:"offerId" bson:"offerId"`
	PromoID      primitive.ObjectID `json:"promoId" bson:"promoId"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	Amount       int                `json:"amount" bson:"amount"`
	Status       string             `json:"status" bson:"status"`
	Transactions []Transaction      `json:"transactions" bson:"transactions"`
}

type Transaction struct {
	Status         string    `json:"status" bson:"status"`
	CreatedAt      time.Time `json:"createdAt" bson:"createdAt"`
	AdditionalInfo string    `json:"additionalInfo" bson:"additionalInfo"`
}
