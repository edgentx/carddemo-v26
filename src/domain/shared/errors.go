package shared

import "errors"

// Common Domain Errors
var (
	ErrUnknownCommand      = errors.New("unknown command")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrInvalidBalance      = errors.New("invalid balance for operation")
	ErrInvariantViolated   = errors.New("domain invariant violated")
)
