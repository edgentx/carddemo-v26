package shared

import "errors"

var (
	// ErrUnknownCommand is returned when a command handler is not found.
	ErrUnknownCommand = errors.New("unknown command")

	// ErrUpstreamNotFound is returned when upstream data sources are missing.
	ErrUpstreamNotFound = errors.New("upstream source not found")

	// ErrInvalidState is returned when an aggregate state prevents an action.
	ErrInvalidState = errors.New("invalid aggregate state")
)
