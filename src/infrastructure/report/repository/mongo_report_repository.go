package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/report/model"
	"github.com/carddemo/project/src/domain/report/repository"
	"github.com/carddemo/project/src/infrastructure/shared"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoReportRepository implements repository.ReportRepository using MongoDB.
type MongoReportRepository struct {
	client *shared.MongoClient
	col    *mongo.Collection
}

// NewMongoReportRepository creates a new MongoReportRepository.
func NewMongoReportRepository(client *shared.MongoClient) *MongoReportRepository {
	return &MongoReportRepository{
		client: client,
		col:    client.Client.Database("carddemo").Collection("reports"),
	}
}

// Get retrieves a report by ID.
func (r *MongoReportRepository) Get(id string) (*model.Report, error) {
	var report model.Report
	ctx, cancel := r.client.Context()
	defer cancel()

	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&report)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &report, nil
}

// Save persists a report aggregate (Upsert).
func (r *MongoReportRepository) Save(report *model.Report) error {
	ctx, cancel := r.client.Context()
	defer cancel()

	filter := bson.M{"_id": report.ID}
	update := bson.M{"$set": report}

	opts := options.Update().SetUpsert(true)
	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

// List retrieves reports with filtering and pagination.
func (r *MongoReportRepository) List(filters repository.FilterOptions, pagination repository.PaginationOptions) ([]*model.Report, string, error) {
	ctx, cancel := r.client.Context()
	defer cancel()

	// Build query
	query := bson.M{}
	if filters.AccountID != nil {
		query["account_id"] = *filters.AccountID
	}
	if filters.Type != nil {
		query["type"] = *filters.Type
	}
	if filters.Status != nil {
		query["status"] = *filters.Status
	}
	if filters.StartDate != nil || filters.EndDate != nil {
		dateQuery := bson.M{}
		if filters.StartDate != nil {
			dateQuery["$gte"] = filters.StartDate
		}
		if filters.EndDate != nil {
			dateQuery["$lte"] = filters.EndDate
		}
		query["created_at"] = dateQuery
	}

	// Handle Cursor Pagination (Sort by created_at desc)
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{Key: "created_at", Value: -1}, {Key: "_id", Value: 1}})

	if pagination.Limit > 0 {
		findOpts.SetLimit(int64(pagination.Limit))
	} else {
		findOpts.SetLimit(10) // Default limit
	}

	if pagination.Cursor != "" {
		decodedCursor, err := decodeCursor(pagination.Cursor)
		if err == nil {
			// Apply cursor filter: find documents created_at < cursor.created_at OR (created_at = cursor.created_at AND id < cursor.id)
			query["$or"] = []bson.M{
				{"created_at": bson.M{"$lt": decodedCursor.CreatedAt}},
				{"created_at": decodedCursor.CreatedAt, "_id": bson.M{"$lt": decodedCursor.ID}},
			}
		}
	}

	cursor, err := r.col.Find(ctx, query, findOpts)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	var results []*model.Report
	if err = cursor.All(ctx, &results); err != nil {
		return nil, "", err
	}

	// Generate next cursor
	nextCursor := ""
	if len(results) > 0 {
		lastReport := results[len(results)-1]
		cursorData := map[string]interface{}{
			"id":         lastReport.ID,
			"created_at": lastReport.CreatedAt,
		}
		jsonData, _ := json.Marshal(cursorData)
		nextCursor = base64.StdEncoding.EncodeToString(jsonData)
	}

	return results, nextCursor, nil
}

// CreateIndexes ensures necessary indexes exist.
func (r *MongoReportRepository) CreateIndexes() error {
	ctx, cancel := r.client.Context()
	defer cancel()

	indexModels := []mongo.IndexModel{
		{Keys: bson.D{{Key: "account_id", Value: 1}}},
		{Keys: bson.D{{Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "account_id", Value: 1}, {Key: "created_at", Value: -1}}},
	}

	_, err := r.col.Indexes().CreateMany(ctx, indexModels)
	return err
}

func decodeCursor(cursor string) (struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}, error) {
	var data struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}

	bytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(bytes, &data)
	return data, err
}
