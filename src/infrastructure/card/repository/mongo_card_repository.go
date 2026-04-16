package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/src/domain/card/repository"
	"github.com/carddemo/project/src/domain/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoCardRepository implements repository.CardRepository.
type mongoCardRepository struct {
	client *mongo.Client
	dbName string
	col    *mongo.Collection
}

// NewMongoCardRepository creates a new MongoDB repository for Cards.
func NewMongoCardRepository(client *mongo.Client, dbName string) repository.CardRepository {
	return &mongoCardRepository{
		client: client,
		dbName: dbName,
		col:    client.Database(dbName).Collection("cards"),
	}
}

// Get retrieves a Card by its aggregate ID.
func (r *mongoCardRepository) Get(ctx context.Context, id string) (*model.Card, error) {
	var doc model.Card
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

// Save persists the state of a Card aggregate.
func (r *mongoCardRepository) Save(ctx context.Context, card *model.Card) error {
	filter := bson.M{"_id": card.ID}
	update := bson.M{"$set": card}

	// Calculate hashed card number for query safety
	hash := sha256.Sum256([]byte(card.CardNumber))
	card.HashedCardNumber = hex.EncodeToString(hash[:])

	opts := options.Update().SetUpsert(true)
	_, err := r.col.UpdateByID(ctx, card.ID, update, opts)
	if err != nil {
		return err
	}

	// TODO: Publish domain events to outbox/event store here.

	return nil
}

// FindByAccountID retrieves all cards associated with a specific account.
func (r *mongoCardRepository) FindByAccountID(ctx context.Context, accountID string) ([]*model.Card, error) {
	filter := bson.M{"account_id": accountID}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cards []*model.Card
	if err = cursor.All(ctx, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}

// FindByNumber retrieves a card by its PAN (card number).
func (r *mongoCardRepository) FindByNumber(ctx context.Context, number string) (*model.Card, error) {
	hash := sha256.Sum256([]byte(number))
	hashStr := hex.EncodeToString(hash[:])

	filter := bson.M{"hashed_card_number": hashStr}
	var doc model.Card
	err := r.col.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, shared.ErrNotFound
		}
		return nil, err
	}
	return &doc, nil
}

// FindByStatusAndType retrieves cards matching a specific status and type.
func (r *mongoCardRepository) FindByStatusAndType(ctx context.Context, status string, cardType string) ([]*model.Card, error) {
	filter := bson.M{"status": status, "type": cardType}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cards []*model.Card
	if err = cursor.All(ctx, &cards); err != nil {
		return nil, err
	}
	return cards, nil
}

// InitializeIndexes creates indexes for the Card collection.
// This should be called on application startup.
func (r *mongoCardRepository) InitializeIndexes(ctx context.Context) error {
	idxModels := []mongo.IndexModel{
		{Keys: bson.D{{Key: "account_id", Value: 1}}},
		{Keys: bson.D{{Key: "hashed_card_number", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "type", Value: 1}}},
	}
	_, err := r.col.Indexes().CreateMany(ctx, idxModels)
	return err
}
