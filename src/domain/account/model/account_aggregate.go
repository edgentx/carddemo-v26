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
	// State fields would be tracked here (e.g., Balance, Status).
	// For OpenAccount, we assume a clean slate or hydration from DB.
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
	case command.OpenAccountCmd:
		return a.handleOpenAccount(c)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleOpenAccount processes the OpenAccountCmd.
func (a *Account) handleOpenAccount(cmd command.OpenAccountCmd) ([]shared.DomainEvent, error) {
	// Scenario: OpenAccountCmd rejected — Account status must be 'Pending' or 'Active'
	if cmd.InitialStatus != "Pending" && cmd.InitialStatus != "Active" {
		return nil, shared.ErrInvalidStatus
	}

	// Create event
	evt := event.NewAccountOpened(a.ID, cmd)
	evt.Payload.AccountID = a.ID
	evt.Payload.UserProfileID = cmd.UserProfileID
	evt.Payload.Status = cmd.InitialStatus
	evt.Payload.AccountType = cmd.AccountType

	// Apply state mutation (optimistic locking + state update)
	// In this green phase, we increment version on successful command.
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
