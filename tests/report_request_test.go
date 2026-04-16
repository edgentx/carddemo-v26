package tests

import (
	"errors"
	"testing"

	"github.com/carddemo/project/src/domain/report/command"
	"github.com/carddemo/project/src/domain/report/model"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReportAggregate_RequestReport_Success tests the successful execution of RequestReportCmd.
func TestReportAggregate_RequestReport_Success(t *testing.T) {
	// Setup: Create a valid Report aggregate.
	reportID := "rep-123"
	agg := model.NewReport(reportID)

	// Build the command.
	cmd := command.RequestReportCmd{
		ReportID:   reportID,
		ConfigID:   "settlement-daily",
		Format:     "csv",
		Parameters: map[string]interface{}{"date": "2023-01-01"},
	}

	// Execute the command.
	events, err := agg.Execute(cmd)

	// Assertions.
	require.NoError(t, err, "Execute should not return an error")
	require.Len(t, events, 1, "Exactly one event should be emitted")

	// Validate Event Type.
	evt := events[0]
	assert.Equal(t, "com.carddemo.report.requested", evt.Type)

	// Validate Event Data structure (it should contain our specific payload).
	reqData, ok := evt.Data.(map[string]interface{})
	require.True(t, ok, "Event data should be a map")
	assert.Equal(t, reportID, reqData["report_id"])
	assert.Equal(t, "settlement-daily", reqData["config_id"])
	assert.Equal(t, "csv", reqData["format"])
}

// TestReportAggregate_RequestReport_Fail_Settle mentNotFinalized tests rejection when data settlement is not finalized.
func TestReportAggregate_RequestReport_Fail_SettlementNotFinalized(t *testing.T) {
	// Setup: Create a Report aggregate in a state that violates the invariant.
	// We assume the aggregate has a way to reach this invalid state for the sake of the scenario.
	// Since we are in RED phase, we might not have the field yet, but we simulate the expectation.
	reportID := "rep-456"
	agg := model.NewReport(reportID)

	// NOTE: In the implementation phase, we would likely need to hydrate this aggregate
	// with a state indicating 'SourceDataNotFinalized'. For now, we pass the command.

	cmd := command.RequestReportCmd{
		ReportID: reportID,
		ConfigID: "settlement-daily",
		Format:   "csv",
	}

	// Execute.
	_, err := agg.Execute(cmd)

	// Assertions.
	// We expect a specific domain error.
	// Since we haven't defined the constant in the domain package yet (it's part of implementation),
	// we check for the error message string or a generic error if the type isn't defined yet.
	// However, the story says "Then the command is rejected with a domain error".
	require.Error(t, err)
	assert.True(t, errors.Is(err, shared.ErrInvalidState), "Should return ErrInvalidState")
}

// TestReportAggregate_RequestReport_Fail_AlreadyArchived tests rejection when report is already archived/immutable.
func TestReportAggregate_RequestReport_Fail_AlreadyArchived(t *testing.T) {
	// Setup: Create a Report aggregate representing an already generated/archived report.
	reportID := "rep-789"
	agg := model.NewReport(reportID)

	cmd := command.RequestReportCmd{
		ReportID: reportID,
		ConfigID: "settlement-daily",
		Format:   "csv",
	}

	// Execute.
	_, err := agg.Execute(cmd)

	// Assertions.
	require.Error(t, err, "Should reject request on immutable aggregate")
	assert.True(t, errors.Is(err, shared.ErrImmutable), "Should return ErrImmutable")
}
