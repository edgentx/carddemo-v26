package model

import (
	"github.com/carddemo/project/src/domain/shared"
)

// Account represents the Account aggregate.
type Account struct {
	shared.AggregateRoot
	ID string
}

// NewAccount creates a new Account instance.
func NewAccount(id string) *Account {
	return &Account{ID: id}
}

// Execute handles commands for the Account aggregate.
func (a *Account) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch cmd.(type) {
	// Command handlers will be implemented here.
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// ID returns the aggregate ID.
func (a *Account) GetID() string {
	return a.ID
}
