package model

import (
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/exportjob/command"
	"github.com/carddemo/project/src/domain/exportjob/event"
	"github.com/carddemo/project/src/domain/shared"
)

// ExportJob represents the ExportJob aggregate.
type ExportJob struct {
	shared.AggregateRoot
	ID string
}

// NewExportJob creates a new ExportJob instance.
func NewExportJob(id string) *ExportJob {
	return &ExportJob{
		ID: id,
		AggregateRoot: shared.AggregateRoot{},
	}
}

// Execute handles commands for the ExportJob aggregate.
func (e *ExportJob) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch c := cmd.(type) {
	case command.InitiateExportCmd:
		return e.handleInitiateExport(c)
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

	// Apply the event to the aggregate (emit it)
	// Using the helper from AggregateRoot to package the event
	domainEvent := shared.NewDomainEvent(
		"com.carddemo.export.initiated", // CloudEvents type convention
		e.ID,
		payload,
	)

	return []shared.DomainEvent{domainEvent}, nil
}

// ID returns the aggregate ID.
func (e *ExportJob) GetID() string {
	return e.ID
}

// ID satisfies the shared.Aggregate interface.
func (e *ExportJob) ID() string {
	return e.ID
}

// ApplyEvent updates the aggregate state based on a domain event.
// Note: In this pure state model, we might not have internal state fields to update yet,
// but the pattern requires processing events if reconstituting from state.
func (e *ExportJob) ApplyEvent(evt shared.DomainEvent) error {
	switch evt.Type {
	case event.EventExportInitiated, "com.carddemo.export.initiated":
		// If we had status fields like 'Status', we would update them here.
		// e.Status = "Initiated"
		return nil
	default:
		return errors.New("unknown event type")
	}
}
