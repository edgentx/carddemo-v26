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
	// Scenario: UpdateAccountStatusCmd rejected — Account status must be 'Pending' or 'Active' to process financial transactions
	// Interpretation: This rule effectively defines the valid set of statuses for the system.
	// If the command attempts to set a status outside this set (e.g. "Frozen", "Locked"), it is rejected.
	validStatuses := map[string]bool{
		"Pending":   true,
		"Active":    true,
		"Suspended": true,
		"Closed":    true,
	}

	if !validStatuses[cmd.NewStatus] {
		return nil, shared.ErrInvalidStatus
	}

	// Scenario: UpdateAccountStatusCmd rejected — Account closure is irreversible and requires a zero balance.
	if cmd.NewStatus == "Closed" {
		if a.Balance != 0 {
			return nil, shared.ErrInvariantViolated
		}
	}

	// Scenario: Successfully execute UpdateAccountStatusCmd
	// 1. Capture old state for event
	oldStatus := a.Status

	// 2. Create Event
	evt := event.NewAccountStatusUpdated(a.ID)
	evt.Payload.AccountID = a.ID
	evt.Payload.OldStatus = oldStatus
	evt.Payload.NewStatus = cmd.NewStatus
	evt.Payload.Reason = cmd.Reason

	// 3. Apply state mutations
	a.Status = cmd.NewStatus
	a.Version++

	return []shared.DomainEvent{evt}, nil
}

// ID satisfies the shared.Aggregate interface.
func (a *Account) ID() string {
	return a.ID
}

// GetID returns the aggregate ID.
func (a *Account) GetID() string {
	return a.ID
}
