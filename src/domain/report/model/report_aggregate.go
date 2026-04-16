package model

import (
	"time"

	"github.com/card-demo/project/src/domain/report/command"
	"github.com/card-demo/project/src/domain/report/event"
	"github.com/card-demo/project/src/domain/shared"
)

// ReportAggregate manages the lifecycle of a report.
type ReportAggregate struct {
	shared.AggregateRoot
	ID        string
	Type      string
	Status    string
	Params    map[string]string
	CreatedAt time.Time
}

// NewReportAggregate creates a new ReportAggregate.
func NewReportAggregate(reportType string, params map[string]string) *ReportAggregate {
	return &ReportAggregate{
		ID:        shared.GenerateID(),
		Type:      reportType,
		Status:    "created",
		Params:    params,
		CreatedAt: time.Now(),
	}
}

// Handle processes commands against the aggregate.
func (a *ReportAggregate) Handle(cmd interface{}) error {
	switch c := cmd.(type) {
	case command.RequestReport:
		return a.requestReport(c)
	case command.ArchiveReport:
		return a.archiveReport(c)
	default:
		return nil
	}
}

// requestReport handles the RequestReport command.
func (a *ReportAggregate) requestReport(cmd command.RequestReport) error {
	if a.Status != "created" {
		return nil // Already requested
	}

	e := event.ReportRequestedEvent{
		Meta: event.Meta{
			AggregateID: a.ID,
			OccurredAt:  time.Now(),
		},
		Type:   cmd.Type,
		Params: cmd.Params,
	}

	a.Apply(e)
	a.AddEvent(e)
	return nil
}

// archiveReport handles the ArchiveReport command.
func (a *ReportAggregate) archiveReport(cmd command.ArchiveReport) error {
	if a.Status == "archived" {
		return nil
	}

	e := event.ReportArchivedEvent{
		Meta: event.Meta{
			AggregateID: a.ID,
			OccurredAt:  time.Now(),
		},
	}

	a.Apply(e)
	a.AddEvent(e)
	return nil
}

// Apply changes the aggregate state based on an event.
func (a *ReportAggregate) Apply(evt interface{}) {
	switch e := evt.(type) {
	case event.ReportRequestedEvent:
		a.Status = "pending"
	case event.ReportArchivedEvent:
		a.Status = "archived"
	}
}
