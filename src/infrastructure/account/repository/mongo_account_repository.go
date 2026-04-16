package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/carddemo/project/src/domain/account/model"
	accountrepo "github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoAccountRepository implements accountrepo.AccountRepository using MongoDB.
type MongoAccountRepository struct {
	collection *mongo.Collection
}

// NewMongoAccountRepository creates a new MongoAccountRepository.
// It initializes the collection and ensures indexes exist.
func NewMongoAccountRepository(collection *mongo.Collection) *MongoAccountRepository {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create indexes. For Account, ID is unique by default (_id), but we ensure it.
	// We also index fields commonly used for queries if necessary.
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "_id", Value: 1}},
	})
	if err != nil {
		// Log warning but don't panic in constructor typically, though for strict tests we might
		fmt.Printf("Warning: failed to create index: %v\n", err)
	}

	return &MongoAccountRepository{
		collection: collection,
	}
}

// Get retrieves an account by ID.
func (r *MongoAccountRepository) Get(id string) (*model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc model.AccountDocument
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}

	agg := &model.Account{}
	agg.FromDocument(&doc)
	return agg, nil
}

// Save persists the account aggregate. It implements optimistic locking.
func (r *MongoAccountRepository) Save(aggregate *model.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := aggregate.ToDocument()

	filter := bson.M{"_id": aggregate.ID}
	// Optimistic Lock: only update if version matches
	filter["version"] = aggregate.Version - 1

	update := bson.M{
		"$set": bson.M{
			"status":     doc.Status,
			"updated_at": time.Now(), // Ensure UpdatedAt is refreshed
			"version":    aggregate.Version,
		},
		"$setOnInsert": bson.M{
			"created_at": doc.CreatedAt,
		},
	}

	opts := options.Update().SetUpsert(true)

	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	// Check for optimistic lock failure
	// If MatchedCount is 0 and UpsertedCount is 0, it means the ID was found but version mismatched (or ID not found, but here we assume upsert handles new).
	// Actually, for a new aggregate, version is 1. filter["version"] == 0.
	// If it exists, it won't match. result.MatchedCount == 0.
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return shared.ErrConcurrencyConflict
	}

	// If it was an upsert (new document), UpdateOne sets the ID in doc if needed, but here we supply ID.
	return nil
}

// Delete removes an account by ID.
func (r *MongoAccountRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// List returns all accounts.
func (r *MongoAccountRepository) List() ([]*model.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []model.AccountDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	accounts := make([]*model.Account, len(docs))
	for i, doc := range docs {
		agg := &model.Account{}
		agg.FromDocument(&doc)
		accounts[i] = agg
	}

	return accounts, nil
}

// Ensure interface compliance
var _ accountrepo.AccountRepository = (*MongoAccountRepository)(nil)
