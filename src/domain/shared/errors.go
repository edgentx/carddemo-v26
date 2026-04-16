package shared

import "errors"

var (
	// ErrNotFound is returned when a requested entity does not exist.
	ErrNotFound = errors.New("entity not found")

	// ErrConcurrencyConflict is returned when an optimistic lock check fails.
	ErrConcurrencyConflict = errors.New("concurrency conflict: version mismatch")
)
