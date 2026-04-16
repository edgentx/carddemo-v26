package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/exportjob/command"
	"github.com/carddemo/project/src/domain/exportjob/event"
	"github.com/carddemo/project/src/domain/exportjob/model"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/stretchr/testify/assert"
)

// mockDomainError allows us to simulate specific domain errors for testing invariants
var mockErrUpstreamNotFound = errors.New("upstream source not found")

func TestExportJob_InitiateExportCmd_Success(t *testing.T) {
	// Arrange
	id := "export-123"
	aggregate := model.NewExportJob(id)
	cmd := command.InitiateExportCmd{
		AggregateID:     id,
		TargetDataset:  "transactions",
		FilterParams:   map[string]interface{}{"start_date": "2023-01-01"},
		UpstreamExists: true,
	}

	// Act
	events, err := aggregate.Execute(cmd)

	// Assert
	assert.NoError(t, err, "Execute should not return an error")
	assert.NotNil(t, events, "Events should not be nil")
	assert.Len(t, events, 1, "One event should be emitted")

	// Validate Domain Event content
	evt := events[0]
	assert.Equal(t, event.EventExportInitiated, evt.Type)
	assert.NotNil(t, evt.Data)

	// Validate Event Payload structure (roughly)
	// Ideally we unmarshal into the struct, but for TDD red phase, checking type is often enough
	// if the struct doesn't exist yet. We do basic checks.
}

func TestExportJob_InitiateExportCmd_UpstreamMissing(t *testing.T) {
	// Arrange
	id := "export-456"
	aggregate := model.NewExportJob(id)
	cmd := command.InitiateExportCmd{
		AggregateID:     id,
		TargetDataset:  "transactions",
		UpstreamExists: false, // Simulating the invariant violation
	}

	// Act
	events, err := aggregate.Execute(cmd)

	// Assert
	assert.Error(t, err, "Execute should return an error when upstream is missing")
	assert.Nil(t, events, "No events should be emitted on error")

	// Ensure it's the specific domain error if possible, or just generic error
	// The requirements say: "An export job cannot proceed if it fails to locate..."
	// We check that the message matches or the error type is correct.
	// Since we are in red phase, we might just check for error existence.
}
