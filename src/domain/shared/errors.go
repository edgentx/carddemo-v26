package shared

import "errors"

// Common domain errors.
var (
	// ErrNotFound is returned when an entity cannot be found in the repository.
	ErrNotFound = errors.New("entity not found")
	// ErrConcurrency is returned when an optimistic locking check fails.
	ErrConcurrency = errors.New("concurrency conflict")
)
