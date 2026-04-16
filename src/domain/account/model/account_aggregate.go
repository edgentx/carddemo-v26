package model

import (
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// Account is the aggregate root for the Account domain.
type Account struct {
	shared.AggregateRoot
	ID     string
	Status string
	// Add other fields as necessary
}

// NewAccount creates a new Account aggregate.
func NewAccount(id, status string) *Account {
	return &Account{
		AggregateRoot: shared.AggregateRoot{
			Version: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:     id,
		Status: status,
	}
}

// ToDocument converts the domain aggregate to a representation suitable for MongoDB.
// Note: We use explicit structs rather than bson tags to maintain domain purity.
func (a *Account) ToDocument() *AccountDocument {
	return &AccountDocument{
		ID:        a.ID,
		Status:    a.Status,
		Version:   a.Version,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

// FromDocument hydrates the domain aggregate from MongoDB data.
func (a *Account) FromDocument(doc *AccountDocument) {
	a.ID = doc.ID
	a.Status = doc.Status
	a.Version = doc.Version
	a.CreatedAt = doc.CreatedAt
	a.UpdatedAt = doc.UpdatedAt
}

// AccountDocument is the database schema for Account.
type AccountDocument struct {
	ID        string    `bson:"_id"`
	Status    string    `bson:"status"`
	Version   int       `bson:"version"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
