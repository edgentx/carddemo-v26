package model

import (
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// UserProfile is the aggregate root for the UserProfile domain.
type UserProfile struct {
	shared.AggregateRoot
	ID    string
	Email string
}

// NewUserProfile creates a new UserProfile aggregate.
func NewUserProfile(id, email string) *UserProfile {
	return &UserProfile{
		AggregateRoot: shared.AggregateRoot{
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:    id,
		Email: email,
	}
}

// ToDocument converts the domain aggregate to a representation suitable for MongoDB.
func (u *UserProfile) ToDocument() *UserProfileDocument {
	return &UserProfileDocument{
		ID:        u.ID,
		Email:     u.Email,
		Version:   u.Version,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromDocument hydrates the domain aggregate from MongoDB data.
func (u *UserProfile) FromDocument(doc *UserProfileDocument) {
	u.ID = doc.ID
	u.Email = doc.Email
	u.Version = doc.Version
	u.CreatedAt = doc.CreatedAt
	u.UpdatedAt = doc.UpdatedAt
}

// UserProfileDocument is the database schema for UserProfile.
type UserProfileDocument struct {
	ID        string    `bson:"_id"`
	Email     string    `bson:"email"`
	Version   int       `bson:"version"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
