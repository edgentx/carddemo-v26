package service

import (
	"testing"

	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/event"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAccountAggregate_RedPhase verifies that the core domain logic
// is failing correctly (or passing if implementation exists).

func TestOpenAccount_RedPhase(t *testing.T) {
	// Arrange
	repo := NewMockAccountRepository()
	service := NewAccountDomainService(repo)
	cmd := &command.OpenAccountCmd{
		UserProfileID: "user-123",
		InitialStatus: "Active",
		AccountType:   "Checking",
	}

	// Act
	result, err := service.ExecuteCommand("", cmd)

	// Assert
	require.NoError(t, err, "Opening account should not error")
	require.NotNil(t, result, "Resulting aggregate should not be nil")

	// Verify State
	assert.Equal(t, "Active", result.Status, "Account status should match command")
	assert.Equal(t, "user-123", result.UserProfileID, "UserProfileID should match command")
	assert.NotEmpty(t, result.ID, "Aggregate ID should be generated")

	// Verify Events
	require.Len(t, result.Events, 1, "One event should be recorded")

	evt, ok := result.Events[0].(*event.AccountOpened)
	require.True(t, ok, "Event should be AccountOpened")

	// Verify Event Envelope (CNCF CloudEvents compliance)
	assert.NotEmpty(t, evt.ID, "CloudEvent ID must be populated")
	assert.Equal(t, "1.0", evt.SpecVersion, "SpecVersion must be 1.0")
	assert.Equal(t, "com.carddemo.account.opened", evt.Type, "Type must be formatted correctly")
	assert.NotEmpty(t, evt.Time, "Time must be set")
	assert.Equal(t, "/account", evt.Source, "Source must be /account")

	// Verify Event Payload
	assert.Equal(t, result.ID, evt.Payload.AccountID, "Payload AccountID must match Aggregate ID")
	assert.Equal(t, "user-123", evt.Payload.UserProfileID, "Payload UserProfileID must match command")
	assert.Equal(t, "Active", evt.Payload.Status, "Payload Status must match command")
	assert.Equal(t, "Checking", evt.Payload.AccountType, "Payload AccountType must match command")

	// Verify Persistence
	saved, err := repo.Get(result.ID)
	require.NoError(t, err)
	assert.Equal(t, result.Status, saved.Status)
}

func TestUpdateAccountStatus_RedPhase(t *testing.T) {
	// Setup: Open an account first
	repo := NewMockAccountRepository()
	service := NewAccountDomainService(repo)
	
	openCmd := &command.OpenAccountCmd{UserProfileID: "user-1", InitialStatus: "Active", AccountType: "Savings"}
	account, _ := service.ExecuteCommand("", openCmd)

	// Now test the update
	updateCmd := &command.UpdateAccountStatusCmd{
		NewStatus: "Frozen",
		Reason:    "Suspicious Activity",
	}

	// Act
	updated, err := service.ExecuteCommand(account.ID, updateCmd)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Frozen", updated.Status)

	// Verify Events
	require.Len(t, updated.Events, 1, "One new event should be recorded")
	statusEvt, ok := updated.Events[0].(*event.AccountStatusUpdated)
	require.True(t, ok)
	assert.Equal(t, "Frozen", statusEvt.Payload.NewStatus)
	assert.Equal(t, "Active", statusEvt.Payload.OldStatus)
	assert.Equal(t, "Suspicious Activity", statusEvt.Payload.Reason)
}

func TestAggregateRootFunctionality(t *testing.T) {
	// Test that events follow the DomainEvent interface
	repo := NewMockAccountRepository()
	service := NewAccountDomainService(repo)
	cmd := &command.OpenAccountCmd{UserProfileID: "u", InitialStatus: "A", AccountType: "C"}
	
	agg, _ := service.ExecuteCommand("", cmd)

	// Check interface compliance
	for _, e := range agg.Events {
		// Ensure methods exist and return expected formats
		assert.NotEmpty(t, e.AggregateID(), "AggregateID must be accessible via interface")
		assert.NotEmpty(t, e.Type(), "Type must be accessible via interface")
	}
}

func TestRepositoryErrorHandling(t *testing.T) {
	repo := NewMockAccountRepository()
	service := NewAccountDomainService(repo)

	// Try to get non-existent account
	_, err := service.ExecuteCommand("non-existent-id", &command.UpdateAccountStatusCmd{NewStatus: "X"})
	
	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, shared.ErrAggregateNotFound)
}
