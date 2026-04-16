package tests

import (
	"testing"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/shared"
)

// Test UpdateAccountStatusCmd Success
func TestAccount_UpdateStatus_Success(t *testing.T) {
	// Given a valid Account aggregate
	acc := model.NewAccount("acc-123")
	acc.Status = "Active"
	acc.Balance = 0.0

	cmd := command.UpdateAccountStatusCmd{
		NewStatus: "Suspended",
		Reason:    "Admin request",
	}

	// When the UpdateAccountStatusCmd command is executed
	events, err := acc.Execute(cmd)

	// Then a account.status.updated event is emitted
	if err != nil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	ev, ok := events[0].(*event.AccountStatusUpdated)
	if !ok {
		t.Fatalf("Expected event type AccountStatusUpdated, got %T", events[0])
	}

	if ev.Type() != "com.carddemo.account.status.updated" {
		t.Errorf("Incorrect event type: %s", ev.Type())
	}
}

// Test UpdateAccountStatusCmd Rejected - Invalid Status for Transactions
func TestAccount_UpdateStatus_Rejected_InvalidStatus(t *testing.T) {
	// Given a valid Account aggregate
	// And a status that violates the rule
	acc := model.NewAccount("acc-456")

	cmd := command.UpdateAccountStatusCmd{
		NewStatus: "Frozen", // Assuming Frozen is not allowed for transaction processing
		Reason:    "Security freeze",
	}

	// When
	_, err := acc.Execute(cmd)

	// Then the command is rejected with a domain error
	if err == nil {
		t.Error("Expected domain error for invalid status, but got nil")
	}
}

// Test UpdateAccountStatusCmd Rejected - Closure with Non-Zero Balance
func TestAccount_UpdateStatus_Rejected_NonZeroBalance(t *testing.T) {
	// Given a Account aggregate that violates: Account closure requires a zero balance
	acc := model.NewAccount("acc-789")
	acc.Status = "Active"
	acc.Balance = 100.50 // Non-zero balance

	cmd := command.UpdateAccountStatusCmd{
		NewStatus: "Closed",
		Reason:    "Customer request",
	}

	// When
	_, err := acc.Execute(cmd)

	// Then the command is rejected with a domain error
	if err == nil {
		t.Error("Expected domain error for closing account with funds, but got nil")
	}

	if err != shared.ErrInvariantViolated {
		t.Errorf("Expected ErrInvariantViolated, got %v", err)
	}
}
