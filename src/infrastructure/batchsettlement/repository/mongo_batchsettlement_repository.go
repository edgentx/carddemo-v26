package repository

import (
	"context"
	"errors"
	"time"

	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/carddemo/project/src/domain/batchsettlement/repository"
	"github.com/carddemo/project/src/domain/transaction/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoBatchSettlementRepository implements repository.BatchSettlementRepository using MongoDB.
type MongoBatchSettlementRepository struct {
	client *mongo.Client
	dbName string
}

// NewMongoBatchSettlementRepository creates a new MongoBatchSettlementRepository.
func NewMongoBatchSettlementRepository(client *mongo.Client, dbName string) repository.BatchSettlementRepository {
	return &MongoBatchSettlementRepository{
		client: client,
		dbName: dbName,
	}
}

func (r *MongoBatchSettlementRepository) getCollection() *mongo.Collection {
	return r.client.Database(r.dbName).Collection("batch_settlements")
}

func (r *MongoBatchSettlementRepository) getTransactionCollection() *mongo.Collection {
	return r.client.Database(r.dbName).Collection("transactions")
}

// Get retrieves a batch settlement by its ID.
func (r *MongoBatchSettlementRepository) Get(id string) (*model.BatchSettlement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var batch model.BatchSettlement
	filter := bson.M{"_id": id}

	err := r.getCollection().FindOne(ctx, filter).Decode(&batch)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("batch settlement not found")
		}
		return nil, err
	}

	return &batch, nil
}

// Save persists the aggregate state and domain events.
func (r *MongoBatchSettlementRepository) Save(aggregate *model.BatchSettlement) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": aggregate.ID}
	update := bson.M{
		"$set": bson.M{
			"merchantId":      aggregate.MerchantID,
			"date":            aggregate.Date,
			"status":          aggregate.Status,
			"transactionIds": aggregate.TransactionIds,
			"createdAt":       aggregate.CreatedAt,
			"version":         aggregate.Version + 1,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.getCollection().UpdateOne(ctx, filter, update, opts)
	return err
}

// List retrieves all batch settlements.
func (r *MongoBatchSettlementRepository) List() ([]*model.BatchSettlement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.getCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var batches []*model.BatchSettlement
	if err = cursor.All(ctx, &batches); err != nil {
		return nil, err
	}

	return batches, nil
}

// BulkInsertTransactions efficiently inserts a slice of transactions into the database.
func (r *MongoBatchSettlementRepository) BulkInsertTransactions(transactions []*model.Transaction) error {
	if len(transactions) == 0 {
		return errors.New("cannot insert empty batch")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	docs := make([]interface{}, len(transactions))
	for i, t := range transactions {
		docs[i] = t
	}

	_, err := r.getTransactionCollection().InsertMany(ctx, docs)
	return err
}

// GetSettlementAggregation groups transactions by merchant and date for reporting/settlement logic.
func (r *MongoBatchSettlementRepository) GetSettlementAggregation(startDate, endDate time.Time) ([]*model.SettlementGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// We must match documents within the date range first.
	// Assuming aggregation is based on the transaction timestamp.
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "timestamp", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lte", Value: endDate},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "merchant", Value: "$merchantId"},
				{Key: "date", Value: bson.D{
					{Key: "$dateToString", Value: bson.D{
						{Key: "format", Value: "%Y-%m-%d"},
						{Key: "date", Value: "$timestamp"},
					}},
				}},
			}},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "merchantId", Value: "$_id.merchant"},
			{Key: "date", Value: bson.D{{Key: "$dateFromString", Value: bson.D{
				{Key: "dateString", Value: "$_id.date"},
				{Key: "format", Value: "%Y-%m-%d"},
			}}}},
			{Key: "total", Value: 1},
			{Key: "count", Value: 1},
		}}},
	}

	cursor, err := r.getTransactionCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.SettlementGroup
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
