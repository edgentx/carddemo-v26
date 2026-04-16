package shared

// AggregateRoot provides base functionality for domain aggregates.
type AggregateRoot struct {
	events []interface{}
}

// RecordEvent adds a domain event to the aggregate.
func (a *AggregateRoot) RecordEvent(event interface{}) {
	a.events = append(a.events, event)
}

// GetEvents returns all recorded events.
func (a *AggregateRoot) GetEvents() []interface{} {
	return a.events
}

// ClearEvents clears the recorded events.
func (a *AggregateRoot) ClearEvents() {
	a.events = nil
}
