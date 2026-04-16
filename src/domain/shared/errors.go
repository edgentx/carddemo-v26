package shared

import "errors"

var (
	// ErrUnknownCommand is returned when the aggregate doesn't support the command.
	ErrUnknownCommand = errors.New("unknown command")

	// ErrAmountMustBePositive is returned when the transaction amount is invalid.
	ErrAmountMustBePositive = errors.New("transaction amount must be strictly greater than zero")

	// ErrAccountNotActive is returned when the account is not in a valid state.
	ErrAccountNotActive = errors.New("account must be in 'Active' status to accept debit or credit transactions")
)
