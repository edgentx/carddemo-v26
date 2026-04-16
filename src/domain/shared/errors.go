package shared

import "errors"

// Common domain errors.
var (
	// ErrUnknownCommand is returned when an aggregate does not recognize the command type.
	ErrUnknownCommand = errors.New("unknown command")
	// ErrConcurrencyConflict is returned when the expected version does not match the actual state.
	ErrConcurrencyConflict = errors.New("concurrency conflict")
)
