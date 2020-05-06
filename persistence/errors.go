package persistence

import "errors"

// package errors
var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)
