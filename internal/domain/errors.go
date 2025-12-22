package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInternal     = errors.New("internal error")
	ErrInvalidInput = errors.New("invalid input")
)
