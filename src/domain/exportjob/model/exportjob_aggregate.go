package model

import (
	"time"

	"github.com/card-demo/project/src/domain/exportjob/command"
	"github.com/card-demo/project/src/domain/exportjob/event"
	"github.com/card-demo/project/src/domain/shared"
)

// ExportJobAggregate manages the lifecycle of an export job.
type ExportJobAggregate struct {
	shared.AggregateRoot
	ID        string
	ReportID  string
	Status    string
	FileKey   string
	CreatedAt time.Time
}

// NewExportJobAggregate creates a new ExportJobAggregate.
func NewExportJobAggregate(reportID string) *ExportJobAggregate {
	return &ExportJobAggregate{
		ID:        shared.GenerateID(),
		ReportID:  reportID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
}

// Handle processes commands against the aggregate.
func (a *ExportJobAggregate) Handle(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.InitiateExport:
		return a.initiateExport(c)
	case command.CompleteExport:
		return a.completeExport(c)
	default:
		return nil
	}
}

// initiateExport handles the InitiateExport command.
func (a *ExportJobAggregate) initiateExport(cmd command.InitiateExport) error {
	if a.Status != "pending" {
		return nil // Idempotency or invalid state check
	}

	e := event.ExportJobInitiatedEvent{
		Meta: event.Meta{
			AggregateID: a.ID,
			OccurredAt:  time.Now(),
		},
		ReportID: cmd.ReportID,
	}

	a.Apply(e)
	a.AddEvent(e)
	return nil
}

// completeExport handles the CompleteExport command.
func (a *ExportJobAggregate) completeExport(cmd command.CompleteExport) error {
	if a.Status != "processing" {
		return nil // Cannot complete if not processing
	}

	e := event.ExportJobCompletedEvent{
		Meta: event.Meta{
			AggregateID: a.ID,
			OccurredAt:  time.Now(),
		},
		FileKey: cmd.FileKey,
	}

	a.Apply(e)
	a.AddEvent(e)
	return nil
}

// Apply changes the aggregate state based on an event.
func (a *ExportJobAggregate) Apply(evt interface{}) {
	switch e := evt.(type) {
	case event.ExportJobInitiatedEvent:
		a.Status = "processing"
	case event.ExportJobCompletedEvent:
		a.Status = "completed"
		a.FileKey = e.FileKey
	}
}
