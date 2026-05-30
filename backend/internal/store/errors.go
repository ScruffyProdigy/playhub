package store

import "errors"

var (
	ErrNotFound       = errors.New("store: not found")
	ErrAlreadyExists  = errors.New("store: already exists")
	ErrAlreadyMatched = errors.New("store: already matched for this game")
)
