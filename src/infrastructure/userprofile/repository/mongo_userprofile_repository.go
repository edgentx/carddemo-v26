package repository

import (
	"context"

	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoUserProfileRepository implements repository.UserProfileRepository.
type MongoUserProfileRepository struct {
	coll *mongo.Collection
}

// NewMongoUserProfileRepository creates a new MongoDB repository.
func NewMongoUserProfileRepository(db *mongo.Database) repository.UserProfileRepository {
	return &MongoUserProfileRepository{
		coll: db.Collection("user_profiles"),
	}
}

// Get retrieves a user profile by ID.
func (r *MongoUserProfileRepository) Get(id string) (*model.UserProfile, error) {
	var agg model.UserProfile
	// Implement mapping from BSON to Aggregate here
	// For brevity, assuming direct mapping or simplified logic
	filter := bson.M{"_id": id}
	err := r.coll.FindOne(context.Background(), filter).Decode(&agg)
	if err != nil {
		return nil, err
	}
	return &agg, nil
}

// Save persists the aggregate state.
// Implementing Atomic Save: Aggregate State + Event Store + Outbox in one txn would go here.
// For this scope, we save the state.
func (r *MongoUserProfileRepository) Save(aggregate *model.UserProfile) error {
	filter := bson.M{"_id": aggregate.ID}
	update := bson.M{"$set": aggregate}

	opts := &mongo.UpdateOptions{}
	// Implement Optimistic Locking with Version check here if needed

	_, err := r.coll.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a user profile.
func (r *MongoUserProfileRepository) Delete(id string) error {
	filter := bson.M{"_id": id}
	_, err := r.coll.DeleteOne(context.Background(), filter)
	return err
}

// List returns all user profiles.
func (r *MongoUserProfileRepository) List() ([]*model.UserProfile, error) {
	cursor, err := r.coll.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var results []*model.UserProfile
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}
