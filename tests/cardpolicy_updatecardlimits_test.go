package tests

import (
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/cardpolicy/command"
	"github.com/carddemo/project/src/domain/cardpolicy/model"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCardPolicy_UpdateCardLimits_Success verifies the happy path for limit updates.
func TestCardPolicy_UpdateCardLimits_Success(t *testing.T) {
	// Setup
	policy := model.NewCardPolicy("policy-123")

	cmd := command.UpdateCardLimitsCmd{
		CardID:             "card-abc",
		DailyLimit:         5000,
		MonthlyLimit:       15000,
		AuthorizationToken: "valid-token",
	}

	// Act
	events, err := policy.Execute(cmd)

	// Assertions
	require.NoError(t, err, "Execute should not return an error")
	assert.Len(t, events, 1, "Exactly one event should be emitted")

	evt, ok := events[0].(model.CardLimitsUpdatedEvent)
	require.True(t, ok, "Event should be of type CardLimitsUpdatedEvent")

	// Verify Event Content
	assert.Equal(t, "com.carddemo.cardpolicy.limits.updated", evt.Type())
	assert.NotEmpty(t, evt.OccurredAt(), "OccurredAt must be set")
	assert.Equal(t, "card-abc", evt.Payload.CardID)
	assert.Equal(t, int64(5000), evt.Payload.DailyLimit)
	assert.Equal(t, int64(15000), evt.Payload.MonthlyLimit)
}

// TestCardPolicy_UpdateCardLimits_InvariantViolation verifies that invalid limits are rejected.
func TestCardPolicy_UpdateCardLimits_InvariantViolation(t *testing.T) {
	// Setup
	policy := model.NewCardPolicy("policy-123")

	// 'InvalidRiskTier' is the specific trigger defined in the aggregate for invariant failure.
	// We pass this as a metadata hint or structured value depending on command design.
	// Assuming we can pass a specific Metadata map or rely on value checks.
	// For this test, we assume the command allows setting a ProfileID.
	cmd := command.UpdateCardLimitsCmd{
		CardID:             "card-xyz",
		DailyLimit:         999999,
		MonthlyLimit:       999999,
		AuthorizationToken: "valid-token",
		ProfileID:          "InvalidRiskTier",
	}

	// Act
	events, err := policy.Execute(cmd)

	// Assertions
	require.Error(t, err, "Execute should return an error for invariant violation")
	assert.Nil(t, events, "No events should be emitted on failure")
	assert.Equal(t, shared.ErrInvariantViolated, err, "Specific invariant error should be returned")
}
