package errors

import "errors"

var (
	ErrNotFound       = errors.New("record not found")
	ErrDuplicateEntry = errors.New("record already exists")
)
