package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidTransaction is returned when a transaction cannot be created.
	ErrInvalidTransaction = errors.New("invalid transaction data")
)

// Transaction represents the domain aggregate for financial transactions.
type Transaction struct {
	ID        string
	CardID    string
	Amount    float64
	Currency  string
	MerchantID string
	Status    string // e.g., "pending", "completed", "reversed"
	Timestamp time.Time
	// Version is used for optimistic locking.
	Version int
}

// NewTransaction creates a new Transaction aggregate.
func NewTransaction(id, cardID string, amount float64, currency, merchantID, status string) (*Transaction, error) {
	if id == "" {
		id = uuid.New().String()
	}
	if amount <= 0 {
		return nil, ErrInvalidTransaction
	}
	return &Transaction{
		ID:        id,
		CardID:    cardID,
		Amount:    amount,
		Currency:  currency,
		MerchantID: merchantID,
		Status:    status,
		Timestamp: time.Now(),
		Version:   0,
	}, nil
}
