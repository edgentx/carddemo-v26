package repository

import (
	"context"
	"errors"

	"github.com/carddemo/project/src/domain/exportjob/model"
	"github.com/carddemo/project/src/domain/exportjob/repository"
	"github.com/carddemo/project/src/infrastructure/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoExportJobRepository implements repository.ExportJobRepository using MongoDB.
type MongoExportJobRepository struct {
	client *shared.MongoClient
	col    *mongo.Collection
}

// NewMongoExportJobRepository creates a new MongoExportJobRepository.
func NewMongoExportJobRepository(client *shared.MongoClient) *MongoExportJobRepository {
	return &MongoExportJobRepository{
		client: client,
		col:    client.Client.Database("carddemo").Collection("export_jobs"),
	}
}

// Get retrieves an export job by ID.
func (r *MongoExportJobRepository) Get(id string) (*model.ExportJob, error) {
	var job model.ExportJob
	ctx, cancel := r.client.Context()
	defer cancel()

	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&job)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

// Save persists an export job aggregate (Upsert).
func (r *MongoExportJobRepository) Save(job *model.ExportJob) error {
	ctx, cancel := r.client.Context()
	defer cancel()

	filter := bson.M{"_id": job.ID}
	update := bson.M{"$set": job}

	opts := options.Update().SetUpsert(true)
	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

// UpdateStatus atomically updates the status of an export job.
func (r *MongoExportJobRepository) UpdateStatus(id string, newStatus string) error {
	ctx, cancel := r.client.Context()
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": newStatus}}

	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

// IncrementRetry atomically increments the retry count.
func (r *MongoExportJobRepository) IncrementRetry(id string) error {
	ctx, cancel := r.client.Context()
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"retry_count": 1}}

	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

// CreateIndexes ensures necessary indexes exist.
func (r *MongoExportJobRepository) CreateIndexes() error {
	ctx, cancel := r.client.Context()
	defer cancel()

	indexModels := []mongo.IndexModel{
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "updated_at", Value: -1}}},
	}

	_, err := r.col.Indexes().CreateMany(ctx, indexModels)
	return err
}
