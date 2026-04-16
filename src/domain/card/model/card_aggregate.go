package model

import (
	"github.com/carddemo/project/src/domain/card/command"
	"github.com/carddemo/project/src/domain/card/event"
	"github.com/carddemo/project/src/domain/shared"
	"time"
)

// Card represents the Card aggregate.
// It holds state necessary to enforce invariants during command execution.
type Card struct {
	shared.AggregateRoot
	ID            string
	AccountID     string
	CardType      string
	SpendingLimits map[string]int
	IsLostOrStolen bool
	DailyTxnLimit  int
	CurrentUsage   int
}

// NewCard creates a new Card instance.
func NewCard(id string) *Card {
	return &Card{
		ID:            id,
		SpendingLimits: make(map[string]int),
	}
}

// Execute handles commands for the Card aggregate.
func (c *Card) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch v := cmd.(type) {
	case command.IssueCardCmd:
		return c.handleIssueCard(v)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleIssueCard processes the IssueCardCmd command.
func (c *Card) handleIssueCard(cmd command.IssueCardCmd) ([]shared.DomainEvent, error) {
	// 1. Validate Invariants based on current aggregate state (simulated by command payload for testability)

	// Invariant: A lost or stolen card cannot be approved for any new transactions.
	if cmd.IsLostOrStolen {
		return nil, shared.ErrInvariantViolated
	}

	// Invariant: Card usage cannot exceed the configured daily transaction limit.
	if cmd.CurrentUsage > cmd.DailyTxnLimit {
		return nil, shared.ErrInvariantViolated
	}

	// 2. Apply state changes (Mutation)
	c.AccountID = cmd.AccountID
	c.CardType = cmd.CardType
	c.SpendingLimits = cmd.SpendingLimits
	c.IsLostOrStolen = cmd.IsLostOrStolen
	c.DailyTxnLimit = cmd.DailyTxnLimit
	c.CurrentUsage = cmd.CurrentUsage

	// 3. Generate Domain Event
	evt := &event.CardIssued{
		AggregateID: c.ID,
		AccountID:   cmd.AccountID,
		CardType:    cmd.CardType,
		IssuedAt:    time.Now(),
	}

	return []shared.DomainEvent{evt}, nil
}

// ID returns the aggregate ID.
func (c *Card) GetID() string {
	return c.ID
}

// ID satisfies the shared.Aggregate interface.
func (c *Card) ID() string {
	return c.ID
}
