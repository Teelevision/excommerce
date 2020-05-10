package persistence

import "errors"

// package errors
var (
	ErrNotFound       = errors.New("not found")
	ErrConflict       = errors.New("conflict")
	ErrNotOwnedByUser = errors.New("not owned by user")
	ErrDeleted        = errors.New("deleted")
)
