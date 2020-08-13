package model

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserConflict = errors.New("user conflict")
)
