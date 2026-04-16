package model

import (
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/exportjob/command"
	"github.com/carddemo/project/src/domain/exportjob/event"
	"github.com/carddemo/project/src/domain/shared"
)

// Status constants to model aggregate state transitions
type ExportStatus string

const (
	StatusInitiated ExportStatus = "initiated"
	StatusCompleted ExportStatus = "completed"
)

// ExportJob represents the ExportJob aggregate.
type ExportJob struct {
	shared.AggregateRoot
	ID     string
	Status ExportStatus
}

// NewExportJob creates a new ExportJob instance.
func NewExportJob(id string) *ExportJob {
	return &ExportJob{
		ID:            id,
		Status:        "", // Initial state
		AggregateRoot: shared.AggregateRoot{},
	}
}

// NewExportJobFromHistory creates an ExportJob from existing events (for testing setup).
// In a real app, this would be done by the Repository loading state.
func (e *ExportJob) ApplyEvent(evt shared.DomainEvent) error {
	switch evt.Type {
	case event.EventExportInitiated, "com.carddemo.export.initiated":
		e.Status = StatusInitiated
	case event.EventExportCompleted, "com.carddemo.export.completed":
		e.Status = StatusCompleted
	default:
		return errors.New("unknown event type")
	}
	return nil
}

// Execute handles commands for the ExportJob aggregate.
func (e *ExportJob) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch c := cmd.(type) {
	case command.InitiateExportCmd:
		return e.handleInitiateExport(c)
	case command.CompleteExportCmd:
		return e.handleCompleteExport(c)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleInitiateExport processes the InitiateExportCmd.
func (e *ExportJob) handleInitiateExport(cmd command.InitiateExportCmd) ([]shared.DomainEvent, error) {
	// Enforce Invariant: An export job cannot proceed if it fails to locate the required upstream source files or data streams.
	// The command provides UpstreamExists to simulate this check statefully or via a injected validation.
	if !cmd.UpstreamExists {
		return nil, shared.ErrUpstreamNotFound
	}

	// Create the domain event
	payload := event.ExportInitiated{
		JobID:         e.ID,
		TargetDataset: cmd.TargetDataset,
		FilterParams:  cmd.FilterParams,
		Timestamp:     time.Now().Unix(),
	}

	// Apply logic locally to ensure consistency before emitting
	e.Status = StatusInitiated

	domainEvent := shared.NewDomainEvent(
		"com.carddemo.export.initiated",
		e.ID,
		payload,
	)

	return []shared.DomainEvent{domainEvent}, nil
}

// handleCompleteExport processes the CompleteExportCmd.
// MARKER: This is the method under test. It starts empty to ensure Red Phase.
func (e *ExportJob) handleCompleteExport(cmd command.CompleteExportCmd) ([]shared.DomainEvent, error) {
	// TODO: Implement state validation
	// TODO: Implement upstream check
	// TODO: Emit event
	return nil, nil
}

// ID satisfies the shared.Aggregate interface.
func (e *ExportJob) ID() string {
	return e.ID
}
