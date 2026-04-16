package model

import (
	"github.com/carddemo/project/src/domain/shared"
)

// BatchSettlement represents the BatchSettlement aggregate.
type BatchSettlement struct {
	shared.AggregateRoot
	ID string
}

// NewBatchSettlement creates a new BatchSettlement instance.
func NewBatchSettlement(id string) *BatchSettlement {
	return &BatchSettlement{ID: id}
}

// Execute handles commands for the BatchSettlement aggregate.
func (b *BatchSettlement) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch cmd.(type) {
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// ID returns the aggregate ID.
func (b *BatchSettlement) GetID() string {
	return b.ID
}

// ID satisfies the shared.Aggregate interface.
func (b *BatchSettlement) ID() string {
	return b.ID
}
