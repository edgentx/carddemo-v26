package shared

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient wraps the mongo.Client to provide application-specific context management.
type MongoClient struct {
	Client *mongo.Client
}

// NewMongoClient creates a new MongoDB client wrapper.
// In a real scenario, this would connect using options from config.
func NewMongoClient(uri string) (*MongoClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &MongoClient{Client: client}, nil
}

// Context returns a context with timeout for database operations.
func (m *MongoClient) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
