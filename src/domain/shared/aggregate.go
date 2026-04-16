package shared

import (
	"github.com/google/uuid"
)

// AggregateRoot provides base functionality for aggregates.
type AggregateRoot struct {
	Events []interface{}
	Version int
}

// AddEvent adds a domain event to the aggregate.
func (a *AggregateRoot) AddEvent(event interface{}) {
	a.Events = append(a.Events, event)
}

// GenerateID creates a new UUID.
func GenerateID() string {
	return uuid.New().String()
}
