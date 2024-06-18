package entity

import "errors"

var (
	ErrUserExists    = errors.New("USER_EXISTS")
	ErrUserNotFound  = errors.New("USER_NOT_FOUND")
	ErrWrongPassword = errors.New("PASSWORD_WRONG")
	ErrInvalidToken  = errors.New("TOKEN_INVALID")
)
