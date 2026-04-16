package shared

// DomainEvent represents a fact that has happened in the domain.
type DomainEvent interface {
	// Type returns the unique name of the event (e.g. com.carddemo.account.opened).
	Type() string
	// AggregateID returns the ID of the aggregate that raised this event.
	AggregateID() string
}
