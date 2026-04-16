package repository

import (
	"context"
	"time"

	"github.com/carddemo/project/src/domain/cardpolicy/model"
)

// CardPolicyRepository defines the storage interface for CardPolicy aggregates.
type CardPolicyRepository interface {
	// Get retrieves a CardPolicy by its aggregate ID.
	Get(ctx context.Context, id string) (*model.CardPolicy, error)

	// Save persists the state of a CardPolicy aggregate.
	Save(ctx context.Context, policy *model.CardPolicy) error

	// FindActivePolicy finds the policy applicable to a card type on a specific date.
	// This typically involves querying by cardType and ensuring the date falls within
	// the effective start and end dates.
	FindActivePolicy(ctx context.Context, cardType string, date time.Time) (*model.CardPolicy, error)
}
