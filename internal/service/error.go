package service

import "errors"

var (
	ErrCourseContentNotFound = errors.New("course has no modules")
)
