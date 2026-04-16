package tests

import (
	"errors"
	"testing"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/domain/shared"
)

// TestAccountOpen_Feature implements the TDD Red Phase for the OpenAccountCmd.
func TestAccountOpen_Feature(t *testing.T) {

	// Scenario: Successfully execute OpenAccountCmd
	t.Run("Success: emit account.opened event", func(t *testing.T) {
		// Given: a valid Account aggregate (empty state for creation)
		repo := mocks.NewMockAccountRepository()
		agg := model.NewAccount("acc_123")

		// And: valid command data
		cmd := command.OpenAccountCmd{
			UserProfileID: "user_01",
			InitialStatus: "Pending",
			AccountType:   "Checking",
		}

		// When: The command is executed
		events, err := agg.Execute(cmd)

		// Then: No error is returned
		if err != nil {
			t.Fatalf("expected success, got error: %v", err)
		}

		// And: Exactly one event is emitted
		if len(events) != 1 {
			t.Fatalf("expected 1 event, got %d", len(events))
		}

		// And: The event is of type AccountOpened (we check type name as the struct isn't exported from package domain yet)
		// Note: In a real implementation, we would type-assert to *event.AccountOpened
		expectedType := "com.carddemo.account.opened"
		if events[0].Type() != expectedType {
			t.Fatalf("expected event type '%s', got '%s'", expectedType, events[0].Type())
		}
	})

	// Scenario: OpenAccountCmd rejected — Account status must be 'Pending' or 'Active'
	t.Run("Failure: invalid initial status rejected", func(t *testing.T) {
		// Given: a valid Account aggregate
		repo := mocks.NewMockAccountRepository()
		agg := model.NewAccount("acc_456")

		// And: An invalid status (e.g. "Closed" or "Random")
		cmd := command.OpenAccountCmd{
			UserProfileID: "user_01",
			InitialStatus: "Closed", // Invalid for opening
			AccountType:   "Checking",
		}

		// When: The command is executed
		events, err := agg.Execute(cmd)

		// Then: A domain error is returned
		if err == nil {
			t.Fatal("expected error for invalid status, got nil")
		}

		// And: The error matches the specific domain invariant
		// We check for a generic error message or a specific wrapped error depending on implementation
		if !errors.Is(err, shared.ErrInvalidStatus) && !errors.Is(err, shared.ErrInvariantViolated) {
			t.Fatalf("expected specific domain error, got: %v", err)
		}

		// And: No events are emitted
		if len(events) != 0 {
			t.Fatalf("expected 0 events on failure, got %d", len(events))
		}
	})

	// Scenario: OpenAccountCmd rejected — Account closure is irreversible
	// Note: The AC phrasing is tricky. It says "Account closure is irreversible...".
	// In the context of OPENING an account, if the status passed implies a closure (Closed),
	// or if the command attempts to open an account that is already closed (via aggregate state logic).
	// Based on the Critical Feedback: "Remove the second if-block checking for ErrInvalidClosure",
	// this test verifies the CORRECT behavior: Open command should fail if status is Closed.
	t.Run("Failure: cannot open account with Closed status", func(t *testing.T) {
		// Given: a valid Account aggregate
		agg := model.NewAccount("acc_789")

		// And: A command with status "Closed"
		cmd := command.OpenAccountCmd{
			UserProfileID: "user_01",
			InitialStatus: "Closed",
			AccountType:   "Checking",
		}

		// When: The command is executed
		_, err := agg.Execute(cmd)

		// Then: The command is rejected (ideally with ErrInvalidStatus)
		if err == nil {
			t.Error("expected error when opening account as Closed, got nil")
		}
	})
}
