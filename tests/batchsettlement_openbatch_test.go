package tests

import (
	"testing"

	"github.com/carddemo/project/src/domain/batchsettlement/command"
	"github.com/carddemo/project/src/domain/batchsettlement/event"
	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOpenBatchCmd_Success verifies that a valid command results in a BatchOpened event.
func TestOpenBatchCmd_Success(t *testing.T) {
	// Given
	agg := model.NewBatchSettlement("batch-123")
	cmd := command.OpenBatchCmd{
		SettlementDate:    "2023-10-27",
		OperationalRegion: "US-EAST",
	}

	// When
	events, err := agg.Execute(cmd)

	// Then
	require.NoError(t, err)
	require.Len(t, events, 1, "Expected exactly one event to be emitted")

	evt, ok := events[0].(event.BatchOpened)
	require.True(t, ok, "Expected event to be of type event.BatchOpened")

	assert.Equal(t, "batch-123", evt.AggregateID())
	assert.Equal(t, "2023-10-27", evt.SettlementDate)
	assert.Equal(t, "US-EAST", evt.OperationalRegion)
	assert.Equal(t, "com.carddemo.batch.opened", evt.Type())
}

// TestOpenBatchCmd_Rejected_UncommittedTransactions verifies the invariant preventing batch opening with pending transactions.
func TestOpenBatchCmd_Rejected_UncommittedTransactions(t *testing.T) {
	// Given
	agg := model.NewBatchSettlement("batch-456")
	// FIX: Explicitly set the state to violate the invariant as per Lead Feedback.
	// The default NewBatchSettlement has HasUncommittedTxns=false.
	agg.HasUncommittedTxns = true

	cmd := command.OpenBatchCmd{
		SettlementDate:    "2023-10-27",
		OperationalRegion: "US-WEST",
	}

	// When
	events, err := agg.Execute(cmd)

	// Then
	require.Error(t, err)
	assert.Nil(t, events, "No events should be emitted when command is rejected")
	assert.Equal(t, model.ErrUncommittedTransactions, err)
}

// TestOpenBatchCmd_Rejected_BalanceMismatch verifies the invariant enforcing financial balance.
func TestOpenBatchCmd_Rejected_BalanceMismatch(t *testing.T) {
	// Given
	agg := model.NewBatchSettlement("batch-789")
	// FIX: Explicitly set the state to violate the invariant as per Lead Feedback.
	// The default NewBatchSettlement has IsBalanced=true.
	agg.IsBalanced = false

	cmd := command.OpenBatchCmd{
		SettlementDate:    "2023-10-27",
		OperationalRegion: "EU-CENTRAL",
	}

	// When
	events, err := agg.Execute(cmd)

	// Then
	require.Error(t, err)
	assert.Nil(t, events, "No events should be emitted when command is rejected")
	assert.Equal(t, model.ErrBalanceMismatch, err)
}
