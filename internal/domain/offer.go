package domain

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PaymentProviderFondy = "fondy"
)

var (
	ErrPaymentProviderNotUsed = errors.New("payment provider is disabled for current offer")
	ErrUnknownPaymentProvider = errors.New("payment provider is not supported")
)

type Offer struct {
	ID            primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name          string               `json:"name" bson:"name"`
	Description   string               `json:"description" bson:"description,omitempty"`
	Benefits      []string             `json:"benefits" bson:"benefits,omitempty"`
	SchoolID      primitive.ObjectID   `json:"schoolId" bson:"schoolId"`
	PackageIDs    []primitive.ObjectID `json:"packages" bson:"packages,omitempty"`
	Price         Price                `json:"price" bson:"price"`
	PaymentMethod PaymentMethod        `json:"paymentMethod" bson:"paymentMethod"`
}

type Price struct {
	Value    uint   `json:"value" bson:"value"`
	Currency string `json:"currency" bson:"currency"`
}

type PaymentMethod struct {
	UsesProvider bool   `json:"usesProvider" bson:"usesProvider"`
	Provider     string `json:"provider" bson:"provider,omitempty"`
}

func (pm PaymentMethod) Validate() error {
	switch pm.Provider {
	case PaymentProviderFondy:
		return nil
	default:
		return errors.New("unknown payment provider")
	}
}
