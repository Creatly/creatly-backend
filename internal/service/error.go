package service

import "errors"

var (
	ErrCourseContentNotFound = errors.New("course has no modules")
	ErrModuleIsNotAvailable  = errors.New("module's content is not available")
	ErrPromocodeExpired      = errors.New("promocode has expired")
)
