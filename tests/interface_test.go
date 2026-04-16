package tests

import (
	"testing"

	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/shared"
)

// TestAggregateContract verifies that the Account aggregate correctly implements the shared.Aggregate interface.
// This is a representative test; in a real setup, you might loop through all aggregates.
func TestAggregateContract(t *testing.T) {
	agg := model.NewAccount("test-id")

	// Check if it satisfies the interface at runtime
	var _ shared.Aggregate = agg

	// Verify ID() method exists and works
	if agg.ID() != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", agg.ID())
	}

	// Verify GetVersion() method exists
	// New aggregates start at version 0 in this base implementation pattern
	if agg.GetVersion() != 0 {
		t.Errorf("Expected version 0, got %d", agg.GetVersion())
	}

	// Verify Execute() signature matches expected behavior
	// (Actual behavior tested in aggregate_test.go)
}

// TestDomainErrorsExist checks that mandatory error variables are defined.
func TestDomainErrorsExist(t *testing.T) {
	if shared.ErrUnknownCommand == nil {
		t.Error("shared.ErrUnknownCommand should not be nil")
	}
	if shared.ErrConcurrencyConflict == nil {
		t.Error("shared.ErrConcurrencyConflict should not be nil")
	}
}
