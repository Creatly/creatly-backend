package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromoCode struct {
	ID                 primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	SchoolId           primitive.ObjectID   `json:"schoolId" bson:"schoolId"`
	Code               string               `json:"code" bson:"code"`
	DiscountPercentage int                  `json:"discountPercentage" bson:"discountPercentage"`
	ExpiresAt          time.Time            `json:"expiresAt" bson:"expiresAt"`
	OfferIDs           []primitive.ObjectID `json:"offerIds" bson:"offerIds"`
}
