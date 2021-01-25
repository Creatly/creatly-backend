package service

import "errors"

var (
	ErrCourseContentNotFound = errors.New("course has no modules")
	ErrModuleIsNotAvailable  = errors.New("module's content is not available")
)
