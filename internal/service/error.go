package service

import "errors"

var (
	ErrUserNotFound         = errors.New("user doesn't exists")
	ErrModuleIsNotAvailable = errors.New("module's content is not available")
	ErrPromocodeExpired     = errors.New("promocode has expired")
	ErrTransactionInvalid   = errors.New("transaction is invalid")
)
