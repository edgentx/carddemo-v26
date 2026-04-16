package tests

import (
	"errors"
	"testing"

	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/transaction/command"
	"github.com/carddemo/project/src/domain/transaction/event"
	"github.com/carddemo/project/src/domain/transaction/model"
)

func TestTransactionReverseCmd(t *testing.T) {
	tests := []struct {
		name           string
		cmd            command.ReverseTransactionCmd
		setupAggregate func() *model.Transaction
		wantErr        error
		validateEvents func(t *testing.T, events []shared.DomainEvent)
	}{
		{
			name: "Successfully execute ReverseTransactionCmd",
			cmd: command.ReverseTransactionCmd{
				TransactionID: "txn-123",
				Amount:        100.0,
				AccountStatus: "Active",
				Reason:        "User requested reversal",
			},
			setupAggregate: func() *model.Transaction {
				return model.NewTransaction("txn-123")
			},
			wantErr: nil,
			validateEvents: func(t *testing.T, events []shared.DomainEvent) {
				if len(events) != 1 {
					t.Fatalf("expected 1 event, got %d", len(events))
				}
				revEvent, ok := events[0].(event.TransactionReversed)
				if !ok {
					t.Fatalf("expected TransactionReversed event, got %T", events[0])
				}
				if revEvent.TransactionID != "txn-123" {
					t.Errorf("expected TransactionID txn-123, got %s", revEvent.TransactionID)
				}
				if revEvent.Reason != "User requested reversal" {
					t.Errorf("expected Reason 'User requested reversal', got %s", revEvent.Reason)
				}
			},
		},
		{
			name: "ReverseTransactionCmd rejected - Transaction amount must be strictly greater than zero",
			cmd: command.ReverseTransactionCmd{
				TransactionID: "txn-456",
				Amount:        0.0,
				AccountStatus: "Active",
				Reason:        "Fraud",
			},
			setupAggregate: func() *model.Transaction {
				return model.NewTransaction("txn-456")
			},
			wantErr: shared.ErrAmountMustBePositive,
			validateEvents: func(t *testing.T, events []shared.DomainEvent) {
				if len(events) != 0 {
					t.Errorf("expected 0 events on error, got %d", len(events))
				}
			},
		},
		{
			name: "ReverseTransactionCmd rejected - Account must be in 'Active' status to accept debit or credit transactions",
			cmd: command.ReverseTransactionCmd{
				TransactionID: "txn-789",
				Amount:        50.0,
				AccountStatus: "Frozen",
				Reason:        "Duplicate",
			},
			setupAggregate: func() *model.Transaction {
				return model.NewTransaction("txn-789")
			},
			wantErr: shared.ErrAccountNotActive,
			validateEvents: func(t *testing.T, events []shared.DomainEvent) {
				if len(events) != 0 {
					t.Errorf("expected 0 events on error, got %d", len(events))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			agg := tt.setupAggregate()

			// Execute
			events, err := agg.Execute(tt.cmd)

			// Validate Error
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			// Validate Events
			if tt.validateEvents != nil {
				tt.validateEvents(t, events)
			}
		})
	}
}
