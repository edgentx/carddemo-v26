package repository

import (
	"github.com/carddemo/project/src/domain/account/model"
)

// AccountRepository defines the storage contract for the Account aggregate.
type AccountRepository interface {
	// Get retrieves an aggregate by ID. Returns ErrAggregateNotFound if not found.
	Get(id string) (*model.AccountAggregate, error)

	// Save persists the aggregate state atomically.
	Save(agg *model.AccountAggregate) error

	// Delete removes an aggregate.
	Delete(id string) error

	// List returns all aggregates.
	List() ([]*model.AccountAggregate, error)
}
