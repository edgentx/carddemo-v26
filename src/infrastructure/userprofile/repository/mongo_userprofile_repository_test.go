package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoCollectionProfile(t *testing.T) *mongo.Collection {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)

	collection := client.Database("carddemo_test").Collection("profiles_test")

	// Clean up
	_, err = collection.DeleteMany(ctx, bson.D{})
	require.NoError(t, err)

	return collection
}

func TestMongoUserProfileRepository_ImplementsInterface(t *testing.T) {
	var _ repository.UserProfileRepository = (*MongoUserProfileRepository)(nil)
}

func TestMongoUserProfileRepository_UniqueIndexes(t *testing.T) {
	collection := getMongoCollectionProfile(t)
	ctx := context.Background()

	// The repository constructor is expected to create indexes.
	// We verify that the indexes exist on the collection.
	// Note: ListIndexes returns the default _id index plus custom ones.

	repo := NewMongoUserProfileRepository(collection)

	// Allow a brief moment for index creation in async environments (though usually sync in tests)
	time.Sleep(100 * time.Millisecond)

	cursor, err := collection.Indexes().List(ctx)
	require.NoError(t, err)

	var indexes []bson.M
	if err := cursor.All(ctx, &indexes); err != nil {
		t.Fatal(err)
	}

	// We expect at least _id and email_1
	foundEmailIndex := false
	for _, idx := range indexes {
		// MongoDB index names usually look like "email_1"
		if name, ok := idx["name"].(string); ok {
			if name == "email_1" {
				foundEmailIndex = true
				// Check uniqueness
				if unique, ok := idx["unique"].(bool); ok {
					assert.True(t, unique, "Email index must be unique")
				}
			}
		}
	}

	assert.True(t, foundEmailIndex, "Unique index on email was not created")

	// Test uniqueness constraint functionality
	profile1 := model.NewUserProfile("u1", "test@example.com")
	require.NoError(t, repo.Save(profile1))

	profile2 := model.NewUserProfile("u2", "test@example.com") // duplicate email
	err = repo.Save(profile2)

	// Should return a duplicate key error (11000)
	assert.Error(t, err)
}

func TestMongoUserProfileRepository_OptimisticLocking(t *testing.T) {
	collection := getMongoCollectionProfile(t)
	repo := NewMongoUserProfileRepository(collection)

	profile := model.NewUserProfile("u-lock", "lock@example.com")
	require.NoError(t, repo.Save(profile))

	// Simulate concurrent modification
	v1, _ := repo.Get("u-lock")
	v2, _ := repo.Get("u-lock")

	v1.Email = "updated1@example.com"
	require.NoError(t, repo.Save(v1))

	v2.Email = "updated2@example.com"
	err := repo.Save(v2)

	assert.Error(t, err)
	assert.Equal(t, shared.ErrConcurrencyConflict, err)
}
