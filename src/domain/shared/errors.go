package shared

import "errors"

var (
	// ErrUnknownCommand is returned when an unregistered command is executed.
	ErrUnknownCommand = errors.New("unknown command")

	// ErrUpstreamNotFound is returned when the required source files or streams are missing.
	ErrUpstreamNotFound = errors.New("upstream source not found")
)
