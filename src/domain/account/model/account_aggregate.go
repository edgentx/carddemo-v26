package model

import (
	"errors"
	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/shared"
)

// AccountAggregate represents the Account entity.
// It handles commands to modify its state and emits events.
type AccountAggregate struct {
	shared.AggregateRoot
	ID            string
	UserProfileID  string
	Status        string
	AccountType   string
}

// NewAccountAggregate creates a new AccountAggregate.
func NewAccountAggregate(id string) *AccountAggregate {
	return &AccountAggregate{
		AggregateRoot: shared.AggregateRoot{},
		ID:            id,
	}
}

// Handle processes commands. It is the entry point for domain logic.
func (a *AccountAggregate) Handle(cmd interface{}) error {
	switch c := cmd.(type) {
	case *command.OpenAccountCmd:
		return a.openAccount(c)
	case *command.UpdateAccountStatusCmd:
		return a.updateStatus(c)
	default:
		return errors.New("unknown command type")
	}
}

// openAccount handles the creation of a new account.
func (a *AccountAggregate) openAccount(cmd *command.OpenAccountCmd) error {
	// 1. Business Invariants/Validation
	if cmd.UserProfileID == "" {
		return shared.ErrValidation
	}
	if cmd.InitialStatus == "" {
		return shared.ErrValidation
	}
	if cmd.AccountType == "" {
		return shared.ErrValidation
	}

	// 2. Apply State Changes
	a.UserProfileID = cmd.UserProfileID
	a.Status = cmd.InitialStatus
	a.AccountType = cmd.AccountType

	// 3. Generate Domain Event
	evt := event.NewAccountOpened(a.ID, cmd)
	// Hydrate event payload from command and new state
	evt.Payload.AccountID = a.ID
	evt.Payload.UserProfileID = cmd.UserProfileID
	evt.Payload.Status = cmd.InitialStatus
	evt.Payload.AccountType = cmd.AccountType

	a.AddEvent(evt)

	return nil
}

// updateStatus handles the status change of an existing account.
func (a *AccountAggregate) updateStatus(cmd *command.UpdateAccountStatusCmd) error {
	// 1. Business Invariants
	if a.Status == cmd.NewStatus {
		return nil // Idempotent or no-op
	}
	if cmd.NewStatus == "" {
		return shared.ErrValidation
	}

	oldStatus := a.Status

	// 2. Apply State Changes
	a.Status = cmd.NewStatus

	// 3. Generate Domain Event
	evt := event.NewAccountStatusUpdated(a.ID)
	evt.Payload.AccountID = a.ID
	evt.Payload.OldStatus = oldStatus
	evt.Payload.NewStatus = cmd.NewStatus
	evt.Payload.Reason = cmd.Reason

	a.AddEvent(evt)

	return nil
}
