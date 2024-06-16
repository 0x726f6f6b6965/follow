package services

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
	ErrSalt         = errors.New("salt error")
	ErrSetFollow    = errors.New("error set follow relationship")
)
