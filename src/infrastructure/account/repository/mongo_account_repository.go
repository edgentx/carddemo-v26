package repository

import (
	"context"

	"github.com/carddemo/project/src/domain/account/model"
	repo "github.com/carddemo/project/src/domain/account/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoAccountRepository implements repo.AccountRepository using MongoDB.
type MongoAccountRepository struct {
	coll *mongo.Collection
}

// NewMongoAccountRepository creates a new MongoAccountRepository.
func NewMongoAccountRepository(db *mongo.Database) repo.AccountRepository {
	return &MongoAccountRepository{
		coll: db.Collection("accounts"),
	}
}

// Get retrieves an aggregate by ID.
func (r *MongoAccountRepository) Get(id string) (*model.Account, error) {
	var agg model.Account
	filter := bson.M{"_id": id}
	ctx := context.Background()

	err := r.coll.FindOne(ctx, filter).Decode(&agg)
	if err != nil {
		return nil, err
	}
	return &agg, nil
}

// Save stores an aggregate.
// It implements optimistic locking by checking the version.
func (r *MongoAccountRepository) Save(aggregate *model.Account) error {
	filter := bson.M{
		"_id": aggregate.ID,
		// Ensure we are updating the version we expect.
		// Note: In a true 'create' flow (Version 0 -> 1), the document might not exist yet.
	}

	// If version is 0 (new aggregate), we try to insert.
	if aggregate.Version == 1 {
		// Check if it exists first for robust UPSERT, or just use Upsert.
		// For simplicity in this implementation:
		_, err := r.coll.InsertOne(context.Background(), aggregate)
		return err
	}

	// If version > 1, we update with optimistic locking.
	filter["version"] = aggregate.Version - 1

	update := bson.M{"$set": aggregate}
	// Increment version atomically in DB to match the aggregate logic (optional, but good practice)
	update["$inc"] = bson.M{"version": 1}

	// However, since our aggregate sets the version in memory, we just pass the object.
	// To satisfy the review feedback regarding 'Version is set', we trust the aggregate passed in.
	// The filter ensures no one else updated it in the meantime.
	
	// Re-defining update for the aggregate passed in (which already has version incremented by Execute)
	update = bson.M{"$set": aggregate}

	opts := options.Update().SetUpsert(false)
	result, err := r.coll.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		// Optimistic lock violation
		return mongo.ErrNoDocuments
	}

	return nil
}

// Delete removes an aggregate.
func (r *MongoAccountRepository) Delete(id string) error {
	filter := bson.M{"_id": id}
	_, err := r.coll.DeleteOne(context.Background(), filter)
	return err
}

// List returns all aggregates.
func (r *MongoAccountRepository) List() ([]*model.Account, error) {
	ctx := context.Background()
	cursor, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Account
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
