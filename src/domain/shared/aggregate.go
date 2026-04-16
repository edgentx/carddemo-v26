package shared

import "time"

// DomainEvent represents a domain event interface.
type DomainEvent interface {
	OccurredOn() time.Time
}

// AggregateRoot is the base struct for all aggregates.
type AggregateRoot struct {
	ID      string   `bson:"_id,omitempty" json:"id"`
	Version int      `bson:"version" json:"version"`
	Events  []DomainEvent `bson:"-" json:"-"`
}
