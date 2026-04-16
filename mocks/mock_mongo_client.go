package mocks

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockMongoClient simulates a MongoDB client for testing.
type MockMongoClient struct {
	mu          sync.Mutex
	Databases   map[string]*MockDatabase
	ShouldError bool
}

// NewMockMongoClient creates a new mock client.
func NewMockMongoClient() *MockMongoClient {
	return &MockMongoClient{
		Databases: make(map[string]*MockDatabase),
	}
}

// Database gets or creates a database mock.
func (m *MockMongoClient) Database(name string, opts ...*options.DatabaseOptions) *MockDatabase {
	m.mu.Lock()
	defer m.mu.Unlock()

	if db, ok := m.Databases[name]; ok {
		return db
	}

	db := &MockDatabase{
		Client:     m,
		Name:       name,
		Collections: make(map[string]*MockCollection),
	}
	m.Databases[name] = db
	return db
}

// MockDatabase simulates a MongoDB database.
type MockDatabase struct {
	Client      *MockMongoClient
	Name        string
	Collections map[string]*MockCollection
	mu          sync.Mutex
}

// Collection gets or creates a collection mock.
func (d *MockDatabase) Collection(name string, opts ...*options.CollectionOptions) *MockCollection {
	d.mu.Lock()
	defer d.mu.Unlock()

	if col, ok := d.Collections[name]; ok {
		return col
	}

	col := &MockCollection{
		Name:   name,
		DBName: d.Name,
		Data:   make([]bson.D, 0),
	}
	d.Collections[name] = col
	return col
}

// MockCollection simulates a MongoDB collection for testing repository logic.
type MockCollection struct {
	Name   string
	DBName string
	Data   []bson.D // Slice of documents to simulate storage
	mu     sync.Mutex
	// InsertError allows simulating DB errors on writes
	InsertError bool
}

// InsertOne simulates inserting a document.
func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.InsertError {
		return nil, &mongo.WriteException{}
	}

	// In a real mock for TDD, we might parse 'document' to bson.D and append to m.Data
	// For now, success is sufficient to drive interface implementation.
	return &mongo.InsertOneResult{}, nil
}

// FindOne simulates finding a single document. Currently returns an empty cursor.
func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	// Return a result that decodes to nothing (document not found) unless configured.
	// This forces the implementation to handle NotFound cases correctly.
	return mongo.NewSingleResultFromDocument(bson.D{}, nil, nil)
}

// Find simulates finding multiple documents.
func (m *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	// Return empty cursor initially.
	return nil, nil
}

// Indexes returns a mock index view.
func (m *MockCollection) Indexes() mongo.IndexView {
	return mongo.IndexView{}
}
