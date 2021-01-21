package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	OfferID      primitive.ObjectID `json:"offerId" bson:"offerId"`
	PromoID      primitive.ObjectID `json:"promoId" bson:"promoId"`
	Status       string             `json:"status" bson:"status"`
	Transactions []Transaction      `json:"transactions" bson:"transactions"`
}

type Transaction struct {
	Status         string `json:"status" bson:"status"`
	CreatedAt      int64  `json:"createdAt" bson:"createdAt"`
	AdditionalInfo string `json:"additionalInfo" bson:"additionalInfo"`
}
