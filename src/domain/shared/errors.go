package shared

import "errors"

var (
	// ErrUnknownCommand is returned when a command cannot be handled by the aggregate.
	ErrUnknownCommand = errors.New("unknown command")

	// ErrInvariantViolated is returned when a domain invariant is violated.
	ErrInvariantViolated = errors.New("invariant violation")
)
