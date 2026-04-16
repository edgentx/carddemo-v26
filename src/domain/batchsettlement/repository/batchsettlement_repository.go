package repository

import (
	"time"

	"github.com/carddemo/project/src/domain/batchsettlement/model"
	"github.com/carddemo/project/src/domain/transaction/model" // Transaction is part of the BatchSettlement logic
)

// BatchSettlementRepository defines the storage and retrieval operations for the BatchSettlement aggregate.
type BatchSettlementRepository interface {
	// Get retrieves a batch settlement by its ID.
	Get(id string) (*model.BatchSettlement, error)

	// Save persists the aggregate state and domain events.
	Save(aggregate *model.BatchSettlement) error

	// List retrieves all batch settlements.
	List() ([]*model.BatchSettlement, error)

	// BulkInsertTransactions efficiently inserts a slice of transactions into the database.
	// This is used during batch creation to link transactions to the settlement.
	BulkInsertTransactions(transactions []*model.Transaction) error

	// GetSettlementAggregation groups transactions by merchant and date for reporting/settlement logic.
	GetSettlementAggregation(startDate, endDate time.Time) ([]*model.SettlementGroup, error)
}
