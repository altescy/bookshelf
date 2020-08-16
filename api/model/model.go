package model

import "errors"

var (
	ErrBookConflict = errors.New("book conflict")
	ErrBookNotFound = errors.New("book not found")
	ErrMimeNotFound = errors.New("mime not found")
	ErrFileConflict = errors.New("file conflict")
	ErrInvalidExt   = errors.New("invalid ext")
)
