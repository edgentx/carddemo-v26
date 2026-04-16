package model

import (
	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/transaction/command"
	"github.com/carddemo/project/src/domain/transaction/event"
)

// Transaction represents the Transaction aggregate.
type Transaction struct {
	shared.AggregateRoot
	ID string
}

// NewTransaction creates a new Transaction instance.
func NewTransaction(id string) *Transaction {
	return &Transaction{ID: id}
}

// Execute handles commands for the Transaction aggregate.
func (t *Transaction) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch c := cmd.(type) {
	case command.SubmitTransactionCmd:
		return t.handleSubmitTransaction(c)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleSubmitTransaction validates the command and applies the resulting events.
func (t *Transaction) handleSubmitTransaction(cmd command.SubmitTransactionCmd) ([]shared.DomainEvent, error) {
	// Invariant 1: Transaction amount must be strictly greater than zero
	if cmd.Amount <= 0 {
		return nil, shared.ErrAmountMustBePositive
	}

	// Invariant 2: Account must be in 'Active' status to accept transactions
	// Note: In a real application, the aggregate might need to load the Account state
	// or the state is validated/passed via the command by the application layer.
	// Based on the Command DTO structure, we validate the status provided.
	if cmd.AccountStatus != "Active" {
		return nil, shared.ErrAccountNotActive
	}

	// Create the domain event
	evt := event.TransactionSubmitted{
		TransactionID:   cmd.TransactionID,
		AccountID:       cmd.AccountID,
		CardID:          cmd.CardID,
		Amount:          cmd.Amount,
		TransactionType: cmd.TransactionType,
	}

	// Record the event internally (if needed for replay/audit) and return
	t.RecordEvent(evt)

	return []shared.DomainEvent{evt}, nil
}

// ID returns the aggregate ID.
func (t *Transaction) GetID() string {
	return t.ID
}

// ID satisfies the shared.Aggregate interface.
func (t *Transaction) ID() string {
	return t.ID
}
