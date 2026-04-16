package repository

import (
	"time"

	"github.com/carddemo/project/src/domain/transaction/model"
)

// TransactionRepository defines the storage and retrieval operations for the Transaction aggregate.
type TransactionRepository interface {
	// Get retrieves a transaction by its ID.
	Get(id string) (*model.Transaction, error)

	// Save persists the aggregate state, domain events, and outbox messages atomically.
	Save(aggregate *model.Transaction) error

	// List retrieves all transactions.
	List() ([]*model.Transaction, error)

	// CreateIndexes ensures the necessary indexes exist for query performance.
	// Specifically on CardId, Timestamp, and Status.
	CreateIndexes() error

	// FindByCardAndDateRange queries transactions for a specific card within a time window.
	FindByCardAndDateRange(cardId string, start, end time.Time) ([]*model.Transaction, error)

	// FindByStatus retrieves transactions by their current status.
	FindByStatus(status string) ([]*model.Transaction, error)
}
