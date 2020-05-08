package controller

import "errors"

// controller errors
var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)
