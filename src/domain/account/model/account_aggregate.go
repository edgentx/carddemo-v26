package model

import (
	"github.com/carddemo/project/src/domain/shared"
)

// AccountStatus represents the status of an account
type AccountStatus string

const (
	AccountStatusOpen  AccountStatus = "Open"
	AccountStatusClosed AccountStatus = "Closed"
)

// Account represents the Account Aggregate
type Account struct {
	shared.AggregateRoot
	ID      string
	OwnerID string
	Status  AccountStatus
	Balance int // in cents
	Version int
}

// NewAccount creates a new Account aggregate
func NewAccount(id, ownerID string) *Account {
	return &Account{
		AggregateRoot: shared.AggregateRoot{},
		ID:            id,
		OwnerID:       ownerID,
		Status:        AccountStatusOpen,
		Balance:       0,
		Version:       0,
	}
}
