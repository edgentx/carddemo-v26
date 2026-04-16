package model

import (
	"github.com/carddemo/project/src/domain/shared"
)

// UserProfile represents the UserProfile aggregate.
type UserProfile struct {
	shared.AggregateRoot
	ID string
}

// NewUserProfile creates a new UserProfile instance.
func NewUserProfile(id string) *UserProfile {
	return &UserProfile{ID: id}
}

// Execute handles commands for the UserProfile aggregate.
func (u *UserProfile) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch cmd.(type) {
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// ID returns the aggregate ID.
func (u *UserProfile) GetID() string {
	return u.ID
}
