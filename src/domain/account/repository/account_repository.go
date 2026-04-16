package repository

import (
	"github.com/carddemo/project/src/domain/account/model"
)

// AccountRepository defines the contract for account persistence.
// It is owned by the domain layer.
type AccountRepository interface {
	// Get retrieves an account aggregate by ID.
	// Returns ErrNotFound if the account does not exist.
	Get(id string) (*model.Account, error)

	// Save persists the account aggregate state.
	// Handles optimistic concurrency via the aggregate's Version field.
	// Atomic: writes state, event store, and outbox in one transaction (conceptually).
	Save(aggregate *model.Account) error

	// Delete removes an account by ID.
	Delete(id string) error

	// List returns all accounts.
	List() ([]*model.Account, error)
}
