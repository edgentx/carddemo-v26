package shared

import "errors"

// Common Domain Errors
var (
	// ErrUnknownCommand is returned when an unregistered command is executed.
	ErrUnknownCommand = errors.New("command not handled")

	// ErrInvalidStatus indicates a business rule violation regarding state.
	ErrInvalidStatus = errors.New("account status must be 'Pending' or 'Active' to process financial transactions")

	// ErrInvariantViolated is a generic wrapper for business logic violations.
	ErrInvariantViolated = errors.New("business rule violation")
)