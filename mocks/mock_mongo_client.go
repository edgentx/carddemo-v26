package mocks

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockMongoClient mimics mongo.Client for unit testing repositories.
// This is a simplified stub to satisfy interfaces and prevent nil panics.
type MockMongoClient struct {
	Database *MockDatabase
}

func NewMockMongoClient() *MockMongoClient {
	return &MockMongoClient{
		Database: &MockDatabase{
			Collections: make(map[string]*MockCollection),
		},
	}
}

func (m *MockMongoClient) Database(name string, opts ...*options.DatabaseOptions) *MockDatabase {
	return m.Database
}

func (m *MockMongoClient) Disconnect(ctx context.Context) error {
	return nil
}

func (m *MockMongoClient) Ping(ctx context.Context) error {
	return nil
}

// MockDatabase mimics mongo.Database.
type MockDatabase struct {
	Collections map[string]*MockCollection
}

func (d *MockDatabase) Collection(name string, opts ...*options.CollectionOptions) *MockCollection {
	if _, ok := d.Collections[name]; !ok {
		d.Collections[name] = &MockCollection{Name: name, Data: make([]bson.Raw, 0)}
	}
	return d.Collections[name]
}

// MockCollection mimics mongo.Collection.
type MockCollection struct {
	Name string
	Data []bson.Raw
}

func (c *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return &mongo.InsertOneResult{InsertedID: "mock_id"}, nil
}

func (c *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	// Return a cursor with 0 results for now to satisfy interface return types in Red Phase.
	// In a more advanced mock, we would decode 'filter' and return static BSON data.
	return nil, nil
}

func (c *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	// Return an empty decoded result (mock mongo.ErrNoDocuments if needed)
	return mongo.NewSingleResultFromDocument(bson.M{}, nil, nil)
}

func (c *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
}

func (c *MockCollection) Indexes() mongo.IndexModelView {
	// This is a tricky interface to mock perfectly without the real driver internals.
	// We return a mock that satisfies the method calls expected by CreateIndexes.
	return &MockIndexView{}
}

// MockIndexView mimics mongo.IndexView.
type MockIndexView struct{}

func (m *MockIndexView) CreateMany(ctx context.Context, models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return []string{"index_1"}, nil
}

func (m *MockIndexView) CreateOne(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	return "index_1", nil
}

func (m *MockIndexView) ListSpecifications(ctx context.Context, opts ...*options.ListSpecificationsOptions) ([]mongo.IndexSpecification, error) {
	return []mongo.IndexSpecification{}, nil
}
