package tests

import (
	"errors"
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/domain/card/command"
	"github.com/carddemo/project/src/domain/card/event"
	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/src/domain/shared"
)

// TestCardAggregate_IssueCard_Success verifies the happy path.
func TestCardAggregate_IssueCard_Success(t *testing.T) {
	// Setup
	repo := mocks.NewMockCardRepository()
	agg := model.NewCard("card-123")

	// Define Command
	cmd := command.IssueCardCmd{
		AccountID:      "acct-456",
		CardType:       "Virtual",
		SpendingLimits: map[string]int{"USD": 1000},
	}

	// Execute
	events, err := agg.Execute(cmd)

	// Assertions
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if len(events) == 0 {
		t.Fatal("expected events to be emitted, got none")
	}

	// Verify specific event type and data
	ev, ok := events[0].(*event.CardIssued)
	if !ok {
		t.Fatalf("expected event.CardIssued, got %T", events[0])
	}

	if ev.AggregateID != "card-123" {
		t.Errorf("expected aggregate ID card-123, got %s", ev.AggregateID)
	}
	if ev.AccountID != "acct-456" {
		t.Errorf("expected account ID acct-456, got %s", ev.AccountID)
	}

	// Ensure state changes persisted if we saved (optional check for flow)
	// _ = repo.Save(agg)
}

// TestCardAggregate_IssueCard_Rejected_DailyLimit verifies the daily transaction limit invariant.
func TestCardAggregate_IssueCard_Rejected_DailyLimit(t *testing.T) {
	// Setup: Create aggregate representing a limit breach scenario
	repo := mocks.NewMockCardRepository()
	agg := model.NewCard("card-limit-123")

	// Define Command with state that triggers the invariant violation
	cmd := command.IssueCardCmd{
		AccountID:      "acct-789",
		CardType:       "Physical",
		SpendingLimits: map[string]int{"USD": 500},

		// Invariant: Card usage cannot exceed the configured daily transaction limit.
		// We simulate this by saying usage is 101, limit is 100.
		DailyTxnLimit: 100,
		CurrentUsage:  101,
	}

	// Execute
	events, err := agg.Execute(cmd)

	// Assertions
	if err == nil {
		t.Fatal("expected error for daily limit violation, got nil")
	}

	if !errors.Is(err, shared.ErrInvariantViolated) {
		t.Logf("expected ErrInvariantViolated, got: %v", err)
	}

	if events != nil {
		t.Error("expected no events to be emitted on failure, got some")
	}
}

// TestCardAggregate_IssueCard_Rejected_LostOrStolen verifies the lost/stolen card invariant.
func TestCardAggregate_IssueCard_Rejected_LostOrStolen(t *testing.T) {
	// Setup: Create aggregate representing a lost/stolen card
	repo := mocks.NewMockCardRepository()
	agg := model.NewCard("card-lost-123")

	// Define Command with state that triggers the invariant violation
	cmd := command.IssueCardCmd{
		AccountID:      "acct-999",
		CardType:       "Virtual",
		SpendingLimits: map[string]int{},

		// Invariant: A lost or stolen card cannot be approved for any new transactions.
		IsLostOrStolen: true,
	}

	// Execute
	events, err := agg.Execute(cmd)

	// Assertions
	if err == nil {
		t.Fatal("expected error for lost/stolen card violation, got nil")
	}

	if !errors.Is(err, shared.ErrInvariantViolated) {
		t.Logf("expected ErrInvariantViolated, got: %v", err)
	}

	if events != nil {
		t.Error("expected no events to be emitted on failure, got some")
	}
}
