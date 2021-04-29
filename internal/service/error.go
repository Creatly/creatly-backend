package service

import "errors"

var (
	ErrUserNotFound            = errors.New("user doesn't exists")
	ErrOfferNotFound           = errors.New("offer doesn't exists")
	ErrPromoNotFound           = errors.New("promocode doesn't exists")
	ErrModuleIsNotAvailable    = errors.New("module's content is not available")
	ErrPromocodeExpired        = errors.New("promocode has expired")
	ErrTransactionInvalid      = errors.New("transaction is invalid")
	ErrUnknownCallbackType     = errors.New("unknown callback type")
	ErrVerificationCodeInvalid = errors.New("verification code is invalid")
	ErrUserAlreadyExists       = errors.New("user with such email already exists")
)
