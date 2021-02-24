package repository

import "errors"

var (
	ErrUserNotFound  = errors.New("user doesn't exists")
	ErrOfferNotFound = errors.New("offer doesn't exists")
	ErrPromoNotFound = errors.New("promocode doesn't exists")
)
