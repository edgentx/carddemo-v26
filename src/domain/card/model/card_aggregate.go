package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/carddemo/project/src/domain/card/command"
	"github.com/carddemo/project/src/domain/card/event"
	"github.com/carddemo/project/src/domain/shared"
)

var (
	// ErrCardAlreadyLost indicates the card is already lost or stolen.
	ErrCardAlreadyLost = errors.New("card is already reported as lost or stolen")
	// ErrLimitExceeded indicates the daily transaction limit has been reached.
	ErrLimitExceeded = errors.New("card usage cannot exceed the configured daily transaction limit")
)

// Card represents the Card Aggregate.
type Card struct {
	shared.AggregateRoot
	ID           string
	AccountID    string
	CardType     string
	Status       string // Active, Lost, Stolen, Closed
	DailyLimit   int
	DailyUsage   int
	IssuedAt     time.Time
	UpdatedAt    time.Time
}

// Handle executes commands against the Card aggregate.
// It returns an error if the command violates business invariants.
func (c *Card) Handle(cmd interface{}) error {
	switch v := cmd.(type) {
	case *command.ReportCardLostCmd:
		return c.handleReportLost(v)
	default:
		return fmt.Errorf("unknown command: %T", cmd)
	}
}

// handleReportLost processes the ReportCardLostCmd.
func (c *Card) handleReportLost(cmd *command.ReportCardLostCmd) error {
	// 1. Check Invariants (State validation)

	// Check if already lost/stolen
	if c.Status == "LOST" || c.Status == "STOLEN" {
		return ErrCardAlreadyLost
	}

	// Check daily limits
	currentUsage := c.DailyUsage
	if cmd.ForceUsage != nil {
		currentUsage = *cmd.ForceUsage
	}

	// Business rule: usage cannot exceed limit.
	if currentUsage > c.DailyLimit {
		return ErrLimitExceeded
	}

	// 2. Apply State Changes
	c.Status = "LOST"
	c.UpdatedAt = time.Now()

	// 3. Emit Event
	e := &event.CardReportedLost{
		AggregateID: c.ID,
		Reason:      cmd.LossReason,
		ReportedBy:  cmd.ReportedBy,
		OccurredAt:  time.Now(),
	}
	c.RecordEvent(e)

	return nil
}
