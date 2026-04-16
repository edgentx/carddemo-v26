package shared

import "github.com/carddemo/project/src/domain/shared/event"

// AggregateRoot is the base for all aggregates
type AggregateRoot struct {
	events []event.DomainEvent
}

// AddEvent adds a domain event to the aggregate
func (a *AggregateRoot) AddEvent(e event.DomainEvent) {
	a.events = append(a.events, e)
}

// GetEvents retrieves all pending events
func (a *AggregateRoot) GetEvents() []event.DomainEvent {
	return a.events
}

// ClearEvents clears the events list
func (a *AggregateRoot) ClearEvents() {
	a.events = nil
}
