package model

import (
	"github.com/carddemo/project/src/domain/shared"
)

// CardPolicy represents the Card Policy Aggregate
type CardPolicy struct {
	shared.AggregateRoot
	ID         string
	AccountID  string
	DailyLimit int
	WeeklyLimit int
	Version    int
}

// NewCardPolicy creates a new CardPolicy aggregate
func NewCardPolicy(id, accountID string, daily, weekly int) *CardPolicy {
	return &CardPolicy{
		AggregateRoot: shared.AggregateRoot{},
		ID:            id,
		AccountID:     accountID,
		DailyLimit:    daily,
		WeeklyLimit:   weekly,
		Version:       0,
	}
}
