package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Promocode struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id"`
	Code               string               `json:"code" bson:"code"`
	DiscountPercentage int                  `json:"discountPercentage" bson:"discountPercentage"`
	ExpiresAt          time.Time            `json:"expiresAt" bson:"expiresAt"`
	OfferIDs           []primitive.ObjectID `json:"offerIds" bson:"offerIds"`
}
