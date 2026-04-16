package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/card/command"
	"github.com/carddemo/project/src/domain/card/event"
	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/mocks"
)

// setupCard creates a fresh Card aggregate for testing.
func setupCard(id, accountID string, dailyLimit int) *model.Card {
	return &model.Card{
		AggregateRoot: model.AggregateRoot{}, // Assuming embedded type from shared
		ID:            id,
		AccountID:     accountID,
		CardType:      "Virtual",
		Status:        "Active",
		DailyLimit:    dailyLimit,
		DailyUsage:    0,
		IssuedAt:      time.Now(),
	}
}

func TestReportCardLost_Success(t *testing.T) {
	// Scenario: Successfully execute ReportCardLostCmd
	repo := mocks.NewMockCardRepository()

	card := setupCard("card-123", "acct-456", 1000)
	repo.Save(card)

	cmd := &command.ReportCardLostCmd{
		CardID:     "card-123",
		LossReason: "Left at bar",
		ReportedBy: "user-1",
	}

	err := card.Handle(cmd)

	if err != nil {
		t.Fatalf("expected successful execution, got error: %v", err)
	}

	// Check Event
	events := card.GetEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	evt, ok := events[0].(*event.CardReportedLost)
	if !ok {
		t.Fatalf("expected CardReportedLost event, got %T", events[0])
	}

	if evt.AggregateID != "card-123" {
		t.Errorf("expected AggregateID card-123, got %s", evt.AggregateID)
	}

	// Check State (Invariants)
	if card.Status != "LOST" {
		t.Errorf("expected status LOST, got %s", card.Status)
	}
}

func TestReportCardLost_Rejected_LimitExceeded(t *testing.T) {
	// Scenario: ReportCardLostCmd rejected — Card usage cannot exceed the configured daily transaction limit
	repo := mocks.NewMockCardRepository()

	card := setupCard("card-789", "acct-456", 100) // Limit is 100
	card.DailyUsage = 90                           // Valid state, but let's assume the command checks usage relative to limit
	repo.Save(card)

	cmd := &command.ReportCardLostCmd{
		CardID:     "card-789",
		LossReason: "Stolen",
		ReportedBy: "user-1",
		ForceUsage: intPtr(150), // Force usage > Limit (Invariant violation)
	}

	err := card.Handle(cmd)

	if err == nil {
		t.Fatal("expected domain error for limit exceeded, got nil")
	}

	if !errors.Is(err, model.ErrLimitExceeded) {
		t.Errorf("expected ErrLimitExceeded, got %v", err)
	}

	// Verify no events were raised
	if len(card.GetEvents()) != 0 {
		t.Errorf("expected 0 events on rejected command, got %d", len(card.GetEvents()))
	}
}

func TestReportCardLost_Rejected_AlreadyLost(t *testing.T) {
	// Scenario: ReportCardLostCmd rejected — A lost or stolen card cannot be approved for any new transactions
	repo := mocks.NewMockCardRepository()

	card := setupCard("card-999", "acct-456", 100)
	card.Status = "STOLEN" // Already in a bad state
	repo.Save(card)

	cmd := &command.ReportCardLostCmd{
		CardID:     "card-999",
		LossReason: "Also lost",
		ReportedBy: "user-1",
	}

	err := card.Handle(cmd)

	if err == nil {
		t.Fatal("expected domain error for already lost/stolen card, got nil")
	}

	if !errors.Is(err, model.ErrCardAlreadyLost) {
		t.Errorf("expected ErrCardAlreadyLost, got %v", err)
	}

	// Verify no events were raised
	if len(card.GetEvents()) != 0 {
		t.Errorf("expected 0 events on rejected command, got %d", len(card.GetEvents()))
	}
}

func intPtr(i int) *int {
	return &i
}
