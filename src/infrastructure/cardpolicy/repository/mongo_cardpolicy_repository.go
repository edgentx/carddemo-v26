package repository

import (
	"context"
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/cardpolicy/model"
	"github.com/carddemo/project/src/domain/cardpolicy/repository"
	"github.com/carddemo/project/src/domain/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoCardPolicyRepository implements repository.CardPolicyRepository.
type mongoCardPolicyRepository struct {
	client *mongo.Client
	dbName string
	col    *mongo.Collection
}

// NewMongoCardPolicyRepository creates a new MongoDB repository for CardPolicies.
func NewMongoCardPolicyRepository(client *mongo.Client, dbName string) repository.CardPolicyRepository {
	return &mongoCardPolicyRepository{
		client: client,
		dbName: dbName,
		col:    client.Database(dbName).Collection("card_policies"),
	}
}

// Get retrieves a CardPolicy by its aggregate ID.
func (r *mongoCardPolicyRepository) Get(ctx context.Context, id string) (*model.CardPolicy, error) {
	var doc model.CardPolicy
	filter := bson.M{"_id": id}
	err := r.col.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &doc, nil
}

// Save persists the state of a CardPolicy aggregate.
func (r *mongoCardPolicyRepository) Save(ctx context.Context, policy *model.CardPolicy) error {
	filter := bson.M{"_id": policy.ID}
	update := bson.M{"$set": policy}
	opts := options.Update().SetUpsert(true)

	_, err := r.col.UpdateByID(ctx, policy.ID, update, opts)
	if err != nil {
		return err
	}

	// TODO: Publish domain events to outbox/event store here.

	return nil
}

// FindActivePolicy finds the policy applicable to a card type on a specific date.
func (r *mongoCardPolicyRepository) FindActivePolicy(ctx context.Context, cardType string, date time.Time) (*model.CardPolicy, error) {
	// Logic: Find policy for cardType where effective_date <= date AND expiration_date >= date.
	// Sort by effective_date descending to get the most recent policy if multiple overlap.
	filter := bson.M{
		"card_type": cardType,
		"effective_date": bson.M{"$lte": date},
		"expiration_date": bson.M{"$gte": date},
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "effective_date", Value: -1}})

	var doc model.CardPolicy
	err := r.col.FindOne(ctx, filter, opts).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &doc, nil
}

// InitializeIndexes creates indexes for the CardPolicy collection.
func (r *mongoCardPolicyRepository) InitializeIndexes(ctx context.Context) error {
	idxModels := []mongo.IndexModel{
		{Keys: bson.D{{Key: "card_type", Value: 1}}},
		{Keys: bson.D{{Key: "effective_date", Value: 1}, {Key: "expiration_date", Value: 1}}},
	}
	_, err := r.col.Indexes().CreateMany(ctx, idxModels)
	return err
}
