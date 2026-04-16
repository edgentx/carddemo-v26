package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/userprofile/model"
	userprofilerepo "github.com/carddemo/project/src/domain/userprofile/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoUserProfileRepository implements userprofilerepo.UserProfileRepository using MongoDB.
type MongoUserProfileRepository struct {
	collection *mongo.Collection
}

// NewMongoUserProfileRepository creates a new MongoUserProfileRepository.
// It initializes the collection and ensures unique indexes exist (e.g., on Email).
func NewMongoUserProfileRepository(collection *mongo.Collection) *MongoUserProfileRepository {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create unique index on email
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Printf("Warning: failed to create unique index on email: %v\n", err)
	}

	return &MongoUserProfileRepository{
		collection: collection,
	}
}

// Get retrieves a user profile by ID.
func (r *MongoUserProfileRepository) Get(id string) (*model.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc model.UserProfileDocument
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	agg := &model.UserProfile{}
	agg.FromDocument(&doc)
	return agg, nil
}

// Save persists the user profile aggregate. It implements optimistic locking.
func (r *MongoUserProfileRepository) Save(aggregate *model.UserProfile) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := aggregate.ToDocument()

	filter := bson.M{"_id": aggregate.ID}
	// Optimistic Lock: only update if version matches
	filter["version"] = aggregate.Version - 1

	update := bson.M{
		"$set": bson.M{
			"email":      doc.Email,
			"updated_at": time.Now(),
			"version":    aggregate.Version,
		},
		"$setOnInsert": bson.M{
			"created_at": doc.CreatedAt,
		},
	}

	opts := options.Update().SetUpsert(true)

	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		// Handle duplicate key error (11000) for unique email
		if mongo.IsDuplicateKeyError(err) {
			// Return a generic error or specific duplicate key error
			return fmt.Errorf("duplicate email")
		}
		return err
	}

	// Check for optimistic lock failure
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return shared.ErrConcurrencyConflict
	}

	return nil
}

// Delete removes a user profile by ID.
func (r *MongoUserProfileRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List returns all user profiles.
func (r *MongoUserProfileRepository) List() ([]*model.UserProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []model.UserProfileDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	profiles := make([]*model.UserProfile, len(docs))
	for i, doc := range docs {
		agg := &model.UserProfile{}
		agg.FromDocument(&doc)
		profiles[i] = agg
	}

	return profiles, nil
}

// Ensure interface compliance
var _ userprofilerepo.UserProfileRepository = (*MongoUserProfileRepository)(nil)
