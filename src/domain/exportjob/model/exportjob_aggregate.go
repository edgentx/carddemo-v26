package model

import (
	"github.com/carddemo/project/src/domain/shared"
)

// ExportJob represents the ExportJob aggregate.
type ExportJob struct {
	shared.AggregateRoot
	ID string
}

// NewExportJob creates a new ExportJob instance.
func NewExportJob(id string) *ExportJob {
	return &ExportJob{ID: id}
}

// Execute handles commands for the ExportJob aggregate.
func (e *ExportJob) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch cmd.(type) {
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// ID returns the aggregate ID.
func (e *ExportJob) GetID() string {
	return e.ID
}

// ID satisfies the shared.Aggregate interface.
func (e *ExportJob) ID() string {
	return e.ID
}
