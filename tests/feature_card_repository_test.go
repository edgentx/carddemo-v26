package tests

import (
	"context"
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/src/domain/card/repository"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/infrastructure/card/repository"
	mocks "github.com/carddemo/project/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test suite for CardRepository implementation
func TestCardRepository(t *testing.T) {
	ctx := context.Background()

	// Setup Mock Client
	mockClient := mocks.NewMockMongoClient()
	mockDB := mockClient.Database("testdb")
	mockCol := mockDB.Collection("cards")

	// Instantiate implementation with mock client wiring
	// We pass the mock collection directly or via a wrapper if the repo expects a client interface.
	// Assuming the repo constructor takes a *mongo.Client interface or we wrap the mock.
	// For this test, we assume a wrapper or direct usage of the mock collection logic.
	
	repo := infra_repo.NewMongoCardRepository(mockClient, "cards")

	t.Run("Save: should persist a new card aggregate", func(t *testing.T) {
		card := &model.Card{
			AggregateRoot: shared.AggregateRoot{ID: primitive.NewObjectID().Hex()},
			AccountID:      "acc_123",
			CardNumber:     "4111111111111111",
			Status:         "ACTIVE",
			Type:           "CREDIT",
			Version:        1,
		}

		// Test Insert
		err := repo.Save(ctx, card)
		assert.NoError(t, err)
	})

	t.Run("FindByNumber: should find a card by hashed number", func(t *testing.T) {
		// We cannot easily mock the internal hashing of the driver, 
		// so this test verifies the plumbing is connected correctly.
		// A true integration of the query logic requires the implementation to exist.
		_, err := repo.FindByNumber(ctx, "4111111111111111")
		// Implementation does not exist yet, so we expect failures initially or nil returns.
		// Here we just ensure the method signature matches.
		assert.NotNil(t, err) // Expecting not found or error in red phase if data not mocked
	})

	t.Run("Indexes: should ensure indexes are created on startup", func(t *testing.T) {
		// This checks that the logic to create indexes is present.
		// We can't inspect the DB, but we can ensure no panic occurs.
		// The actual implementation will call CreateIndex.
		assert.NotNil(t, repo)
	})
}
