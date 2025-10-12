package ierr

import "errors"

var (
	ErrNotFound     = errors.New("record not found")
	ErrConflict     = errors.New("record already exists or causes a conflict")
	ErrInvalidInput = errors.New("input validation failed")
)
