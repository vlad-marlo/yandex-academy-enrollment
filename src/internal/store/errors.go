package store

import "errors"

var (
	ErrNoContent     = errors.New("")
	ErrDoesNotExists = errors.New("record does not exists")
)
