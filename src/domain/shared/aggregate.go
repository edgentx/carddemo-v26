package shared

// Aggregate defines the contract for domain aggregates.
type Aggregate interface {
	// Execute handles a command and returns resulting domain events.
	Execute(cmd interface{}) ([]DomainEvent, error)
	// ID returns the unique identifier of the aggregate.
	ID() string
	// GetVersion returns the current version of the aggregate for optimistic locking.
	GetVersion() int
}

// AggregateRoot provides a base implementation for state change tracking.
type AggregateRoot struct {
	version int
	events  []DomainEvent
}

// AddEvent appends a domain event to the pending list.
func (a *AggregateRoot) AddEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

// ClearEvents resets the pending events list, usually after persistence.
func (a *AggregateRoot) ClearEvents() {
	a.events = nil
}

// GetEvents returns the currently pending domain events.
func (a *AggregateRoot) GetEvents() []DomainEvent {
	return a.events
}
