package controller

import "errors"

// controller errors
var (
	ErrNotFound  = errors.New("not found")
	ErrConflict  = errors.New("conflict")
	ErrForbidden = errors.New("forbidden")
	ErrDeleted   = errors.New("deleted")
	ErrLocked    = errors.New("locked")
)
