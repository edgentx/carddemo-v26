package tests

import (
	"errors"
	"testing"

	"github.com/carddemo/project/src/domain/batchsettlement/command"
	"github.com/carddemo/project/src/domain/batchsettlement/event"
	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/carddemo/project/src/domain/shared"
)

func TestBatchSettlement_ReconcileBatchCmd(t *testing.T) {
	tests := []struct {
		name             string
		aggregate        *model.BatchSettlement
		cmd              command.ReconcileBatchCmd
		expectedEvents   int
		expectedEventType string
		expectedError    error
	}{
		{
			name: "Scenario: Successfully execute ReconcileBatchCmd",
			aggregate: &model.BatchSettlement{
				AggregateRoot:      shared.AggregateRoot{},
				ID:                 "batch-123",
				HasUncommittedTxns: false,
				IsBalanced:         true,
				TotalDebits:        5000,
				TotalCredits:       5000,
			},
			cmd: command.ReconcileBatchCmd{
				BatchID:             "batch-123",
				ExpectedTotalDebits:  5000,
				ExpectedTotalCredits: 5000,
			},
			expectedEvents:    1,
			expectedEventType: "com.carddemo.batch.reconciled",
			expectedError:     nil,
		},
		{
			name: "Scenario: ReconcileBatchCmd rejected - uncommitted transactions",
			aggregate: &model.BatchSettlement{
				AggregateRoot:      shared.AggregateRoot{},
				ID:                 "batch-456",
				HasUncommittedTxns: true, // Violates invariant
				IsBalanced:         true,
				TotalDebits:        1000,
				TotalCredits:       1000,
			},
			cmd: command.ReconcileBatchCmd{
				BatchID:             "batch-456",
				ExpectedTotalDebits:  1000,
				ExpectedTotalCredits: 1000,
			},
			expectedEvents:    0,
			expectedEventType: "",
			expectedError:     model.ErrUncommittedTransactions,
		},
		{
			name: "Scenario: ReconcileBatchCmd rejected - totals mismatch",
			aggregate: &model.BatchSettlement{
				AggregateRoot:      shared.AggregateRoot{},
				ID:                 "batch-789",
				HasUncommittedTxns: false,
				IsBalanced:         true,
				TotalDebits:        2000, // Actual state
				TotalCredits:       2000,
			},
			cmd: command.ReconcileBatchCmd{
				BatchID:             "batch-789",
				ExpectedTotalDebits:  5000, // Incorrect expectation
				ExpectedTotalCredits: 5000,
			},
			expectedEvents:    0,
			expectedEventType: "",
			expectedError:     model.ErrInvalidReconciliationTotals,
		},
		{
			// Critical Test based on Feedback: Ensures we don't accidentally accept
			// balanced inputs if they don't match the aggregate state.
			name: "Scenario: ReconcileBatchCmd rejected - balanced input but mismatched state",
			aggregate: &model.BatchSettlement{
				AggregateRoot:      shared.AggregateRoot{},
				ID:                 "batch-101",
				HasUncommittedTxns: false,
				IsBalanced:         true,
				TotalDebits:        1000,
				TotalCredits:       1000,
			},
			cmd: command.ReconcileBatchCmd{
				BatchID:             "batch-101",
				ExpectedTotalDebits:  2000,
				ExpectedTotalCredits: 2000, // Inputs are balanced, but WRONG
			},
			expectedEvents:    0,
			expectedEventType: "",
			expectedError:     model.ErrInvalidReconciliationTotals,
		},
		{
			// Critical Test based on Feedback: Debits != Credits mismatch against state
			name: "Scenario: ReconcileBatchCmd rejected - complex mismatch",
			aggregate: &model.BatchSettlement{
				AggregateRoot:      shared.AggregateRoot{},
				ID:                 "batch-202",
				HasUncommittedTxns: false,
				IsBalanced:         true,
				TotalDebits:        1000,
				TotalCredits:       2000,
			},
			cmd: command.ReconcileBatchCmd{
				BatchID:             "batch-202",
				ExpectedTotalDebits:  1000,
				ExpectedTotalCredits: 1000,
			},
			expectedEvents:    0,
			expectedEventType: "",
			expectedError:     model.ErrInvalidReconciliationTotals,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			events, err := tt.aggregate.Execute(tt.cmd)

			// Then - Error Check
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}

			// Then - Event Count Check
			if len(events) != tt.expectedEvents {
				t.Errorf("Expected %d events, got %d", tt.expectedEvents, len(events))
			}

			// Then - Event Content Check (only if success expected)
			if tt.expectedEvents > 0 && len(events) > 0 {
				e := events[0]
				if e.Type() != tt.expectedEventType {
					t.Errorf("Expected event type %s, got %s", tt.expectedEventType, e.Type())
				}

				// Verify payload matches aggregate state (not just cmd input)
				reconciledEvent, ok := e.(event.BatchReconciled)
				if !ok {
					t.Errorf("Expected event.BatchReconciled, got %T", e)
				}

				if reconciledEvent.TotalDebits != tt.aggregate.TotalDebits {
					t.Errorf("Event payload TotalDebits %d does not match aggregate state %d", reconciledEvent.TotalDebits, tt.aggregate.TotalDebits)
				}
				if reconciledEvent.TotalCredits != tt.aggregate.TotalCredits {
					t.Errorf("Event payload TotalCredits %d does not match aggregate state %d", reconciledEvent.TotalCredits, tt.aggregate.TotalCredits)
				}
			}
		})
	}
}
