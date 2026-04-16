package repository

import (
	"github.com/carddemo/project/src/domain/userprofile/model"
)

// UserProfileRepository defines the contract for user profile persistence.
// It is owned by the domain layer.
type UserProfileRepository interface {
	// Get retrieves a user profile aggregate by ID.
	// Returns ErrNotFound if the profile does not exist.
	Get(id string) (*model.UserProfile, error)

	// Save persists the user profile aggregate state.
	// Handles optimistic concurrency via the aggregate's Version field.
	// Atomic: writes state, event store, and outbox in one transaction (conceptually).
	Save(aggregate *model.UserProfile) error

	// Delete removes a user profile by ID.
	Delete(id string) error

	// List returns all user profiles.
	List() ([]*model.UserProfile, error)
}
