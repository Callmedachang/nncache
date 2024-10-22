package nncache

import "errors"

var (
	// ErrEntryNotFound is an error type struct which is returned when entry was not found for provided key
	ErrEntryNotFound = errors.New("Entry not found")
	// ErrEntryIsDead is an error type struct which is returned when entry has past it's life window
	ErrEntryIsDead = errors.New("entry is dead")
)

