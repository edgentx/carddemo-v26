package model

import (
	"github.com/carddemo/project/src/domain/shared"
)

// CardPolicy represents the CardPolicy aggregate.
type CardPolicy struct {
	shared.AggregateRoot
	ID string
}

// NewCardPolicy creates a new CardPolicy instance.
func NewCardPolicy(id string) *CardPolicy {
	return &CardPolicy{ID: id}
}

// Execute handles commands for the CardPolicy aggregate.
func (c *CardPolicy) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch cmd.(type) {
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// ID returns the aggregate ID.
func (c *CardPolicy) GetID() string {
	return c.ID
}
