package service

import "errors"

var (
	ErrModuleIsNotAvailable  = errors.New("module's content is not available")
	ErrPromocodeExpired      = errors.New("promocode has expired")
	ErrTransactionInvalid    = errors.New("transaction is invalid")
)
