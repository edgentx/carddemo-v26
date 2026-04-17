package model

import (
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// UserProfile represents the UserProfile Aggregate.
type UserProfile struct {
	shared.AggregateRoot
	ID        string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUserProfile creates a new UserProfile aggregate.
func NewUserProfile(id, firstName, lastName, email string) (*UserProfile, error) {
	return &UserProfile{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// UpdateDetails updates the profile details.
func (u *UserProfile) UpdateDetails(firstName, lastName, email string) {
	u.FirstName = firstName
	u.LastName = lastName
	u.Email = email
	u.UpdatedAt = time.Now()
}
