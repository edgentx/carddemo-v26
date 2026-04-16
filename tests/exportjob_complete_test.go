package tests

import (
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/exportjob/command"
	"github.com/carddemo/project/src/domain/exportjob/event"
	"github.com/carddemo/project/src/domain/exportjob/model"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompleteExportCmd_Success verifies the happy path for completing an export job.
// This enforces the requirement: 
// "Given a valid ExportJob aggregate... export.completed event is emitted"
func TestCompleteExportCmd_Success(t *testing.T) {
	// Setup: Create a valid ExportJob aggregate
	// We simulate an 'Initiated' state by creating an aggregate and manually setting status
	// or applying an Initiated event. For this Unit Test, we manually set the state
	// to verify the Command Handler logic specifically.
	job := model.NewExportJob("job-123")
	job.Status = model.StatusInitiated

	cmd := command.CompleteExportCmd{
		AggregateID:  "job-123",
		ExportID:     "export-abc",
		RecordCount:  5000,
		ManifestData: "{\"files\":[\"data.csv\"]}",
	}

	// Execute: Run the command against the aggregate
	events, err := job.Execute(cmd)

	// Assertions:
	require.NoError(t, err, "Execute should not return an error")
	require.Len(t, events, 1, "Exactly one event should be emitted")

	// Verify Event Content
	evt := events[0]
	assert.Equal(t, "com.carddemo.export.completed", evt.Type)
	assert.Equal(t, "job-123", evt.AggregateID)

	// Verify Payload
	payload, ok := evt.Data.(event.ExportCompleted)
	require.True(t, ok, "Payload should be ExportCompleted type")
	assert.Equal(t, int64(5000), payload.RecordCount)
	assert.Equal(t, "{\"files\":[\"data.csv\"]}", payload.ManifestData)
	assert.Greater(t, payload.Timestamp, int64(0))

	// Verify Aggregate State Update (if any)
	// We expect the status to move to Completed
	assert.Equal(t, model.StatusCompleted, job.Status)
}

// TestCompleteExportCmd_RejectedIfUpstreamMissing verifies the domain invariant:
// "An export job cannot proceed if it fails to locate the required upstream source files or data streams"
// The feedback implies the Command carries the validation context (or we check against aggregate state).
// Based on InitiateExportCmd pattern which used 'UpstreamExists', we expect similar logic or state check.
// Here we assume the aggregate tracks upstream availability or the command provides the flag.
// Given the 'Scenario' text, we simulate a situation where upstream check fails.
// However, the command definition doesn't have a boolean flag like Initiate.
// The prompt feedback says: 'handleCompleteExport lacks any state checks on the aggregate itself.'
// This implies we should check Aggregate State. If state is 'Initiated', we assume it was valid.
// If we fail to locate files DURING completion (e.g. writing manifest failed), we reject.
// To make this test meaningful and 'Red', we will assume the Command/Handler logic MUST check an internal flag
// or the Command should return an error. We will write the test expecting an error.
func TestCompleteExportCmd_RejectedIfUpstreamMissing(t *testing.T) {
	// Setup: Create an aggregate that represents a failed upstream check
	// For example, a state 'UpstreamFailed'
	job := model.NewExportJob("job-404")
	
	// We use a custom status to represent the 'Upstream Missing' scenario described in the AC
	// "Given a ExportJob aggregate that violates..."
	job.Status = "UpstreamMissing" 

	cmd := command.CompleteExportCmd{
		AggregateID:  "job-404",
		ExportID:     "export-404",
		RecordCount:  0,
		ManifestData: "",
	}

	// Execute
	events, err := job.Execute(cmd)

	// Assertions
	require.Error(t, err, "Execute should return an error when upstream is missing")
	assert.Nil(t, events, "No events should be emitted on failure")
	assert.Equal(t, shared.ErrUpstreamNotFound, err)
}

// TestCompleteExportCmd_InvalidState checks if we try to complete an already completed job.
func TestCompleteExportCmd_InvalidState(t *testing.T) {
	job := model.NewExportJob("job-999")
	job.Status = model.StatusCompleted // Already done

	cmd := command.CompleteExportCmd{
		AggregateID:  "job-999",
		ExportID:     "export-999",
		RecordCount:  100,
		ManifestData: "{}",
	}

	events, err := job.Execute(cmd)

	require.Error(t, err)
	assert.Equal(t, shared.ErrInvalidState, err)
	assert.Nil(t, events)
}

// TestEvents_CloudEventsFormat checks compliance with CloudEvents spec requirements.
func TestEvents_CloudEventsFormat(t *testing.T) {
	// Prepare
	job := model.NewExportJob("job-format")
	job.Status = model.StatusInitiated

	cmd := command.CompleteExportCmd{
		AggregateID:  "job-format",
		ExportID:     "exp",
		RecordCount:  1,
		ManifestData: "data",
	}

	events, _ := job.Execute(cmd)
	evt := events[0]

	// CloudEvents 1.0 Attributes
	assert.Equal(t, "1.0", evt.SpecVersion)
	assert.NotEmpty(t, evt.ID, "Event ID must be a UUID")
	assert.NotEmpty(t, evt.Source, "Source must be the service name")
	assert.Equal(t, "application/json", evt.DataContentType)
}
