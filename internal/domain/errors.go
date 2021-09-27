package domain

import "errors"

var (
	ErrUserNotFound            = errors.New("user doesn't exists")
	ErrVerificationCodeInvalid = errors.New("verification code is invalid")
	ErrOfferNotFound           = errors.New("offer doesn't exists")
	ErrPromoNotFound           = errors.New("promocode doesn't exists")
	ErrCourseNotFound          = errors.New("course not found")
	ErrUserAlreadyExists       = errors.New("user with such email already exists")
	ErrModuleIsNotAvailable    = errors.New("module's content is not available")
	ErrPromocodeExpired        = errors.New("promocode has expired")
	ErrTransactionInvalid      = errors.New("transaction is invalid")
	ErrUnknownCallbackType     = errors.New("unknown callback type")
	ErrSendPulseIsNotConnected = errors.New("sendpulse is not connected")
	ErrStudentBlocked          = errors.New("student is blocked by the admin")
)
