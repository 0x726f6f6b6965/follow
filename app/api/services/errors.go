package services

import "errors"

var (
	ErrTimeout         = errors.New("timeout")
	ErrInvalidInput    = errors.New("invalid input")
	ErrInvalidResponse = errors.New("invalid response")
)
