package services

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrFollow       = errors.New("error set follow relationship")
)
