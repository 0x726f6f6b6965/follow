package services

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrSetFollow    = errors.New("error set follow relationship")
)
