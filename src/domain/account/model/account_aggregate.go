package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/shared"
)

// Account represents the Account aggregate.
type Account struct {
	ID      string
	Version int
	Status  string
	Balance float64
}

// NewAccount creates a new Account instance.
func NewAccount(id string) *Account {
	return &Account{
		ID:      id,
		Version: 0,
	}
}

// Execute handles commands for the Account aggregate.
func (a *Account) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch c := cmd.(type) {
	case command.UpdateAccountStatusCmd:
		return a.handleUpdateAccountStatus(c)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleUpdateAccountStatus processes the UpdateAccountStatusCmd.
func (a *Account) handleUpdateAccountStatus(cmd command.UpdateAccountStatusCmd) ([]shared.DomainEvent, error) {
	// 1. Invariant: Account status must be 'Pending' or 'Active' to process financial transactions
	// Interpretation: If the NEW status is intended for transaction processing, it must be valid.
	// However, the prompt says: "UpdateAccountStatusCmd rejected — Account status must be 'Pending' or 'Active' to process financial transactions"
	// Context implies we are validating the state transition or the resulting state.
	// Let's assume the rule applies to the Target Status if it implies activity, OR the source status.
	// Given standard DDD, let's enforce that the Account is in a valid state (Pending/Active) to accept commands generally.
	// BUT, the AC specifically lists this under UpdateAccountStatusCmd.
	// Let's enforce: You can only move TO Pending or Active (or maybe FROM). 
	// Let's stick to the text: "Account status must be 'Pending' or 'Active' to process financial transactions".
	// This is a bit ambiguous. Let's assume the command tries to set a status that allows transactions.
	// Valid statuses: Pending, Active, Suspended, Closed.
	if cmd.NewStatus != "Pending" && cmd.NewStatus != "Active" && cmd.NewStatus != "Suspended" && cmd.NewStatus != "Closed" {
		return nil, errors.New("invalid target status")
	}

	// 2. Invariant: Account closure is irreversible and requires a zero balance.
	if cmd.NewStatus == "Closed" {
		if a.Balance != 0 {
			return nil, shared.ErrInvariantViolated // "Account closure... requires a zero balance"
		}
	}

	// If we are closing, check invariants.
	// For this red phase, we will just return nil to fail tests until implemented.
	return nil, errors.New("not implemented")
}
