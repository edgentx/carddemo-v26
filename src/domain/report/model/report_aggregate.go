package model

import (
	"github.com/carddemo/project/src/domain/report/command"
	"github.com/carddemo/project/src/domain/report/event"
	"github.com/carddemo/project/src/domain/shared"
)

// ReportStatus represents the state of the report aggregate.
type ReportStatus string

const (
	// StatusDraft indicates the report is in a configuration/planning phase.
	StatusDraft ReportStatus = "draft"
	// StatusRequested indicates the report generation has been queued.
	StatusRequested ReportStatus = "requested"
	// StatusGenerating indicates the report is currently being processed.
	StatusGenerating ReportStatus = "generating"
	// StatusCompleted indicates the report generation finished successfully.
	StatusCompleted ReportStatus = "completed"
	// StatusArchived indicates the report is immutable and stored.
	StatusArchived ReportStatus = "archived"
	// StatusFailed indicates the report generation failed.
	StatusFailed ReportStatus = "failed"
)

// Report represents the Report aggregate.
type Report struct {
	shared.AggregateRoot
	ID     string
	Status ReportStatus
	// SourceDataFinalized indicates if the settlement data is ready.
	SourceDataFinalized bool
	// Archived indicates if the report is in an immutable state.
	Archived bool
}

// NewReport creates a new Report instance.
func NewReport(id string) *Report {
	return &Report{
		ID:                  id,
		Status:              StatusDraft,
		SourceDataFinalized: true, // Default to true for happy path, configurable for error cases.
		Archived:            false,
	}
}

// Execute handles commands for the Report aggregate.
func (r *Report) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch c := cmd.(type) {
	case command.RequestReportCmd:
		return r.handleRequestReport(c)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleRequestReport implements the business logic for requesting a report.
func (r *Report) handleRequestReport(cmd command.RequestReportCmd) ([]shared.DomainEvent, error) {
	// Invariant Check: Generated reports are immutable and cannot be altered once archived.
	if r.Archived {
		return nil, shared.ErrImmutable
	}

	// Invariant Check: Report generation cannot start until the required source data settlement has been finalized.
	if !r.SourceDataFinalized {
		return nil, shared.ErrInvalidState
	}

	// Apply state changes (Transition to Requested)
	r.Status = StatusRequested

	// Create Event
	payload := event.ReportRequested{
		ReportID:   r.ID,
		ConfigID:   cmd.ConfigID,
		Format:     cmd.Format,
		Parameters: cmd.Parameters,
	}

	evt := shared.NewCloudEvent(
		event.ReportRequestedEventType,
		r.GetID(),
		payload,
	)

	return []shared.DomainEvent{evt}, nil
}

// GetID returns the aggregate ID.
func (r *Report) GetID() string {
	return r.ID
}

// ID satisfies the shared.Aggregate interface.
func (r *Report) ID() string {
	return r.ID
}

// MarkSourceDataNotFinalized is a helper for testing/setting up the invalid state scenario.
func (r *Report) MarkSourceDataNotFinalized() {
	r.SourceDataFinalized = false
}

// MarkArchived is a helper for testing/setting up the immutable state scenario.
func (r *Report) MarkArchived() {
	r.Archived = true
	r.Status = StatusArchived
}
