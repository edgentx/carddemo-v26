package model

import (
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// CardStatus represents the state of a card
type CardStatus string

const (
	CardStatusIssued    CardStatus = "Issued"
	CardStatusActive    CardStatus = "Active"
	CardStatusSuspended CardStatus = "Suspended"
	CardStatusBlocked   CardStatus = "Blocked"
	CardStatusClosed    CardStatus = "Closed"
)

// Card represents the Card Aggregate
type Card struct {
	shared.AggregateRoot
	ID             string
	AccountID      string
	CardType       string
	Status         CardStatus
	SpendingLimits map[string]int
	IssuedAt       time.Time
	Version        int
}

// NewCard creates a new Card aggregate
func NewCard(id, accountID, cardType string, limits map[string]int) *Card {
	return &Card{
		AggregateRoot:  shared.AggregateRoot{},
		ID:             id,
		AccountID:      accountID,
		CardType:       cardType,
		Status:         CardStatusIssued,
		SpendingLimits: limits,
		IssuedAt:       time.Now(),
		Version:        0,
	}
}
