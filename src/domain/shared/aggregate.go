package shared

import "time"

// AggregateRoot is the base struct for all aggregates.
type AggregateRoot struct {
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}
