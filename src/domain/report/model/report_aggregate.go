package model

import (
	"github.com/carddemo/project/src/domain/shared"
)

// Report represents the Report aggregate.
type Report struct {
	shared.AggregateRoot
	ID string
}

// NewReport creates a new Report instance.
func NewReport(id string) *Report {
	return &Report{ID: id}
}

// Execute handles commands for the Report aggregate.
func (r *Report) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch cmd.(type) {
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// ID returns the aggregate ID.
func (r *Report) GetID() string {
	return r.ID
}

// ID satisfies the shared.Aggregate interface.
func (r *Report) ID() string {
	return r.ID
}
