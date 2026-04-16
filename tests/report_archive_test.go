package tests

import (
	"testing"

	"github.com/carddemo/project/src/domain/report/command"
	"github.com/carddemo/project/src/domain/report/model"
	"github.com/carddemo/project/src/domain/report/repository"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestArchiveReport_Success verifies the happy path where a report is successfully archived.
func TestArchiveReport_Success(t *testing.T) {
	// Setup
	repo := mocks.NewMockReportRepository()
	agg := model.NewReport("rep-123")
	agg.Status = model.StatusCompleted // Assume generation is done
	
	err := repo.Save(agg)
	require.NoError(t, err)

	// Given
	cmd := command.ArchiveReportCmd{
		ReportID:        "rep-123",
		StorageLocation: "s3://bucket/path/to/report.pdf",
	}

	// When
	reloaded, err := repo.Get("rep-123")
	require.NoError(t, err)

	events, err := reloaded.Execute(cmd)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.Len(t, events, 1)

	// Verify Event Type
	evt := events[0]
	assert.Equal(t, "com.carddemo.report.archived", evt.Type())

	// Verify State Changes (Invariant enforcement: status becomes archived)
	assert.Equal(t, model.StatusArchived, reloaded.Status)
	assert.True(t, reloaded.Archived)
}

// TestArchiveReport_RejectedIfNotFinalized verifies the command fails if source data isn't ready.
// Scenario: Report generation cannot start until the required source data settlement has been finalized.
// Note: The prompt description implies "start", but archiving implies "finish". We interpret this invariant
// as applying to state-changing actions. If SourceDataFinalized is false, we cannot Archive.
func TestArchiveReport_RejectedIfNotFinalized(t *testing.T) {
	// Setup
	repo := mocks.NewMockReportRepository()
	agg := model.NewReport("rep-456")
	agg.Status = model.StatusCompleted
	// Violation: Source data not finalized
	agg.MarkSourceDataNotFinalized() 

	err := repo.Save(agg)
	require.NoError(t, err)

	// Given
	cmd := command.ArchiveReportCmd{
		ReportID:        "rep-456",
		StorageLocation: "s3://bucket/path/to/report.pdf",
	}

	// When
	reloaded, _ := repo.Get("rep-456")
	events, err := reloaded.Execute(cmd)

	// Then
	assert.Error(t, err)
	assert.Nil(t, events)
	assert.ErrorIs(t, err, shared.ErrInvalidState)
}

// TestArchiveReport_RejectedIfImmutable verifies the command fails if the report is already archived.
// Scenario: Generated reports are immutable and cannot be altered once archived.
func TestArchiveReport_RejectedIfImmutable(t *testing.T) {
	// Setup
	repo := mocks.NewMockReportRepository()
	agg := model.NewReport("rep-789")
	// Violation: Report is already archived
	agg.Status = model.StatusArchived
	agg.Archived = true

	err := repo.Save(agg)
	require.NoError(t, err)

	// Given
	cmd := command.ArchiveReportCmd{
		ReportID:        "rep-789",
		StorageLocation: "s3://bucket/path/to/report-v2.pdf", // Attempting to update location
	}

	// When
	reloaded, _ := repo.Get("rep-789")
	events, err := reloaded.Execute(cmd)

	// Then
	assert.Error(t, err)
	assert.Nil(t, events)
	assert.ErrorIs(t, err, shared.ErrImmutable)
}

// TestArchiveReport_UnknownCommand ensures type safety in the Execute method.
func TestReport_Execute_UnknownCommand(t *testing.T) {
	// Setup
	agg := model.NewReport("rep-000")

	// When: Passing a command struct that doesn't match expected types
	events, err := agg.Execute("string command")

	// Then
	assert.Error(t, err)
	assert.Nil(t, events)
	assert.ErrorIs(t, err, shared.ErrUnknownCommand)
}