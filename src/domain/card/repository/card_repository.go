package repository

import (
	"context"

	"github.com/carddemo/project/src/domain/card/model"
)

// CardRepository defines the storage interface for Card aggregates.
type CardRepository interface {
	// Get retrieves a Card by its aggregate ID.
	Get(ctx context.Context, id string) (*model.Card, error)

	// Save persists the state of a Card aggregate.
	Save(ctx context.Context, card *model.Card) error

	// FindByAccountID retrieves all cards associated with a specific account.
	FindByAccountID(ctx context.Context, accountID string) ([]*model.Card, error)

	// FindByNumber retrieves a card by its PAN (card number).
	FindByNumber(ctx context.Context, number string) (*model.Card, error)

	// FindByStatusAndType retrieves cards matching a specific status and type.
	FindByStatusAndType(ctx context.Context, status string, cardType string) ([]*model.Card, error)
}
