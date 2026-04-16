package model

import (
	"github.com/carddemo/project/src/domain/shared"
)

// Card represents the Card aggregate.
type Card struct {
	shared.AggregateRoot
	ID string
}

// NewCard creates a new Card instance.
func NewCard(id string) *Card {
	return &Card{ID: id}
}

// Execute handles commands for the Card aggregate.
func (c *Card) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch cmd.(type) {
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// ID returns the aggregate ID.
func (c *Card) GetID() string {
	return c.ID
}
