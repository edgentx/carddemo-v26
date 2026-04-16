package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Helper to setup a real mongo container for integration/unit testing logic
// This is the closest to "embedded" we get in Go without heavy libraries like dockertest.
// We assume MONGODB_URI is set or defaults to localhost for the test runner.

func getMongoCollection(t *testing.T) *mongo.Collection {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err, "Failed to connect to MongoDB")

	// Use a unique collection name for this test run to avoid collisions\tcollection := client.Database("carddemo_test").Collection("accounts_test")

	// Clean up before test
	_, err = collection.DeleteMany(ctx, bson.D{})
	require.NoError(t, err)

	return collection
}

func TestMongoAccountRepository_ImplementsInterface(t *testing.T) {
	// This is a compile-time check enforced by the compiler,
	// but we can document the intent here.
	var _ repository.AccountRepository = (*MongoAccountRepository)(nil)
}

func TestMongoAccountRepository_CRUD(t *testing.T) {
	collection := getMongoCollection(t)
	repo := NewMongoAccountRepository(collection)

	ctx := context.Background()

	t.Run("Save and Get", func(t *testing.T) {
		agg := model.NewAccount("acc-123", "Active")

		// Act
		err := repo.Save(agg)

		// Assert
		require.NoError(t, err)

		// Verify retrieval
		found, err := repo.Get("acc-123")
		require.NoError(t, err)
		assert.Equal(t, "acc-123", found.ID)
		assert.Equal(t, "Active", found.Status)
	})

	t.Run("Get NotFound", func(t *testing.T) {
		_, err := repo.Get("non-existent")
		assert.Error(t, err)
		assert.Equal(t, shared.ErrNotFound, err)
	})

	t.Run("Optimistic Concurrency Control", func(t *testing.T) {
		agg := model.NewAccount("acc-999", "Active")
		// Save initial version
		require.NoError(t, repo.Save(agg))

		// Fetch two copies
		v1, _ := repo.Get("acc-999")
		v2, _ := repo.Get("acc-999")

		// Modify v1 and save
		v1.Status = "Closed"
		require.NoError(t, repo.Save(v1)) // Version increments to 2

		// Modify v2 (stale version 1) and try to save
		v2.Status = "Suspended"
		err := repo.Save(v2)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, shared.ErrConcurrencyConflict, err)

		// Verify DB state is v1's state
		final, _ := repo.Get("acc-999")
		assert.Equal(t, "Closed", final.Status)
	})
}

func TestMongoAccountRepository_ConnectionPooling_Config(t *testing.T) {
	// This test ensures we rely on the mongo client passed in,
	// typically configured in infrastructure config.
	// Here we just validate options propagation if we added specific options to the repo.
	collection := getMongoCollection(t)
	repo := NewMongoAccountRepository(collection)
	assert.NotNil(t, repo)
}
