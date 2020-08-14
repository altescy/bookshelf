package model

import "errors"

var (
	ErrBookConflict = errors.New("book conflict")
	ErrBookNotFound = errors.New("book not found")
)
