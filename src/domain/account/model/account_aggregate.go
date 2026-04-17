package model

import (
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/shared"
)

// Error definitions
var (
	ErrInvalidStatus = errors.New("invalid status")
	ErrClosed       = errors.New("account is closed")
)

// Account represents the Account Aggregate.
type Account struct {
	shared.AggregateRoot
	ID            string
	UserProfileID string
	Status        string
	AccountType   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewAccount creates a new Account aggregate.
func NewAccount(id, userProfileID, initialStatus, accountType string) (*Account, error) {
	if initialStatus == "" {
		initialStatus = "pending"
	}
	a := &Account{
		ID:            id,
		UserProfileID: userProfileID,
		Status:        initialStatus,
		AccountType:   accountType,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Record Domain Event
	evt := event.NewAccountOpened(id, command.OpenAccountCmd{
		UserProfileID: userProfileID,
		InitialStatus: initialStatus,
		AccountType:   accountType,
	})
	a.RecordEvent(evt)

	return a, nil
}

// Execute handles commands dispatched to the aggregate.
func (a *Account) Execute(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.UpdateAccountStatusCmd:
		return a.updateStatus(c)
	default:
		return shared.ErrUnknownCommand
	}
}

func (a *Account) updateStatus(cmd command.UpdateAccountStatusCmd) error {
	if a.Status == "closed" {
		return ErrClosed
	}

	oldStatus := a.Status
	a.Status = cmd.NewStatus
	a.UpdatedAt = time.Now()

	evt := event.NewAccountStatusUpdated(a.ID)
	evt.Payload.OldStatus = oldStatus
	evt.Payload.NewStatus = cmd.NewStatus
	evt.Payload.Reason = cmd.Reason

	a.RecordEvent(evt)
	return nil
}
