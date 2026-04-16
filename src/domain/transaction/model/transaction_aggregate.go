package model

import (
	"github.com/carddemo/project/src/domain/shared"
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
	switch cmd.(type) {
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// ID returns the aggregate ID.
func (t *Transaction) GetID() string {
	return t.ID
}
