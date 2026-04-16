package tests

import (
	"testing"

	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/transaction/command"
	"github.com/carddemo/project/src/domain/transaction/event"
	"github.com/carddemo/project/src/domain/transaction/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSubmitTransactionCmd_Success tests the happy path where a transaction is valid.
func TestSubmitTransactionCmd_Success(t *testing.T) {
	// Given
	aggregate := model.NewTransaction("txn_123")
	cmd := command.SubmitTransactionCmd{
		TransactionID:   "txn_123",
		AccountID:       "acc_456",
		CardID:          "card_789",
		Amount:          100.50,
		TransactionType: "debit",
		AccountStatus:   "Active",
	}

	// When
	evts, err := aggregate.Execute(cmd)

	// Then
	require.NoError(t, err, "Execute should not return an error")
	require.Len(t, evts, 1, "Exactly one event should be emitted")

	// Verify Event Type and Content
	submittedEvent, ok := evts[0].(event.TransactionSubmitted)
	require.True(t, ok, "Event should be of type TransactionSubmitted")

	assert.Equal(t, "txn_123", submittedEvent.TransactionID)
	assert.Equal(t, 100.50, submittedEvent.Amount)
	assert.Equal(t, "debit", submittedEvent.TransactionType)
}

// TestSubmitTransactionCmd_Rejected_InvalidAmount tests the rejection when amount <= 0.
func TestSubmitTransactionCmd_Rejected_InvalidAmount(t *testing.T) {
	// Given
	aggregate := model.NewTransaction("txn_123")
	// We iterate various invalid amounts to ensure robustness
	invalidAmounts := []float64{0, -10, -0.01}

	for _, amount := range invalidAmounts {
		t.Run("AmountValidation", func(t *testing.T) {
			cmd := command.SubmitTransactionCmd{
				TransactionID:   "txn_123",
				AccountID:       "acc_456",
				CardID:          "card_789",
				Amount:          amount,
				TransactionType: "debit",
				AccountStatus:   "Active", // Account is valid, amount is the issue
			}

			// When
			evts, err := aggregate.Execute(cmd)

			// Then
			require.Error(t, err, "Execute should return an error for invalid amount")
			assert.ErrorIs(t, err, shared.ErrAmountMustBePositive, "Error should be ErrAmountMustBePositive")
			assert.Nil(t, evts, "No events should be emitted on rejection")
		})
	}
}

// TestSubmitTransactionCmd_Rejected_InactiveAccount tests the rejection when account status is not Active.
func TestSubmitTransactionCmd_Rejected_InactiveAccount(t *testing.T) {
	// Given
	aggregate := model.NewTransaction("txn_123")
	// Inactive statuses: 'Inactive', 'Frozen', 'Closed'
	invalidStatuses := []string{"Inactive", "Frozen", "Closed", "Pending"}

	for _, status := range invalidStatuses {
		t.Run("AccountStatusValidation", func(t *testing.T) {
			cmd := command.SubmitTransactionCmd{
				TransactionID:   "txn_123",
				AccountID:       "acc_456",
				CardID:          "card_789",
				Amount:          100.00, // Valid amount
				TransactionType: "credit",
				AccountStatus:   status, // The violation
			}

			// When
			evts, err := aggregate.Execute(cmd)

			// Then
			require.Error(t, err, "Execute should return an error for inactive account")
			assert.ErrorIs(t, err, shared.ErrAccountNotActive, "Error should be ErrAccountNotActive")
			assert.Nil(t, evts, "No events should be emitted on rejection")
		})
	}
}
