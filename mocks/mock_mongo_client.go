package mocks

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MockMongoCollection is a mock for mongo.Collection
type MockMongoCollection struct {
	InsertOneFunc func(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	FindOneFunc   func(ctx context.Context, filter interface{}) *mongo.SingleResult
	UpdateOneFunc func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

func (m *MockMongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	if m.InsertOneFunc != nil {
		return m.InsertOneFunc(ctx, document)
	}
	return &mongo.InsertOneResult{}, nil
}

func (m *MockMongoCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	if m.FindOneFunc != nil {
		return m.FindOneFunc(ctx, filter)
	}
	// Return an empty decoded result error by default
	return mongo.NewSingleResultFromDocument(bson.M{}, nil, nil)
}

func (m *MockMongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if m.UpdateOneFunc != nil {
		return m.UpdateOneFunc(ctx, filter, update, opts...)
	}
	return &mongo.UpdateResult{}, nil
}

// IndexView mock to satisfy collection.Indexes()
type MockIndexView struct {
	CreateOneFunc func(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error)
	ListFunc      func(ctx context.Context, opts ...*options.ListIndexesOptions) (cursor *mongo.Cursor, err error)
}

func (m *MockMongoCollection) Indexes() mongo.IndexView {
	return &MockIndexView{}
}

func (m *MockIndexView) CreateOne(ctx context.Context, model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	if m.CreateOneFunc != nil {
		return m.CreateOneFunc(ctx, model, opts...)
	}
	return "idx_test", nil
}

func (m *MockIndexView) List(ctx context.Context, opts ...*options.ListIndexesOptions) (*mongo.Cursor, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, opts...)
	}
	return nil, nil
}
