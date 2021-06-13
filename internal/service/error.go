package service

import "errors"

var (
	ErrModuleIsNotAvailable    = errors.New("module's content is not available")
	ErrPromocodeExpired        = errors.New("promocode has expired")
	ErrTransactionInvalid      = errors.New("transaction is invalid")
	ErrUnknownCallbackType     = errors.New("unknown callback type")
	ErrVerificationCodeInvalid = errors.New("verification code is invalid")
)
