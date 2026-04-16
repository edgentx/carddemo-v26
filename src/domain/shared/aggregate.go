package shared

// Aggregate is the interface that all aggregates must implement.
type Aggregate interface {
	ID() string
	Execute(cmd interface{}) ([]DomainEvent, error)
}
