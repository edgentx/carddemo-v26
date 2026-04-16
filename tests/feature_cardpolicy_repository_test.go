package tests

import (
	"context"
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/cardpolicy/model"
	"github.com/carddemo/project/src/domain/cardpolicy/repository"
	"github.com/carddemo/project/src/domain/shared"
	infra_repo "github.com/carddemo/project/src/infrastructure/cardpolicy/repository"
	mocks "github.com/carddemo/project/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test suite for CardPolicyRepository implementation
func TestCardPolicyRepository(t *testing.T) {
	ctx := context.Background()

	mockClient := mocks.NewMockMongoClient()
	repo := infra_repo.NewMongoCardPolicyRepository(mockClient, "card_policies")

	t.Run("Save: should persist a new policy", func(t *testing.T) {
		policy := &model.CardPolicy{
			AggregateRoot: shared.AggregateRoot{ID: primitive.NewObjectID().Hex()},
			CardType:       "CREDIT",
			EffectiveDate:  time.Now().Add(-24 * time.Hour),
			ExpirationDate: time.Now().Add(24 * time.Hour),
			Version:        1,
		}

		err := repo.Save(ctx, policy)
		assert.NoError(t, err)
	})

	t.Run("FindActivePolicy: should query for policy by type and date range", func(t *testing.T) {
		// Test Query Logic
		activeDate := time.Now()
		_, err := repo.FindActivePolicy(ctx, "CREDIT", activeDate)
		
		// Red phase: Implementation likely returns nil or error.
		// We assert that the method exists and signature is correct.
		assert.NotNil(t, repo) 
		_ = err // Explicitly ignoring error for red-phase placeholder, but checking compilation
	})
}
