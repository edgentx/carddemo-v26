package shared

import "errors"

// Common Domain Errors
var (
	// ErrUnknownCommand is returned when a command type is not recognized by the aggregate.
	ErrUnknownCommand = errors.New("unknown command")

	// ErrInvalidState is returned when a business rule invariant is violated based on state.
	ErrInvalidState = errors.New("invalid state: operation not allowed")

	// ErrImmutable is returned when trying to modify an immutable entity or aggregate.
	ErrImmutable = errors.New("entity is immutable and cannot be modified")
)
