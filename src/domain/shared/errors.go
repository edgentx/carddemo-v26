package shared

import "errors"

var (
	// ErrUnknownCommand is returned when an unregistered command is executed.
	ErrUnknownCommand = errors.New("command not recognized")

	// ErrInvariantViolated is returned when a business rule is broken.
	ErrInvariantViolated = errors.New("business rule violation")
)
