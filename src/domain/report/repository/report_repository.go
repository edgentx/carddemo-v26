package repository

import (
	"time"

	"github.com/carddemo/project/src/domain/report/model"
)

// FilterOptions holds the criteria for querying reports.
type FilterOptions struct {
	AccountID *string
	Type      *string
	Status    *string
	StartDate *time.Time
	EndDate   *time.Time
}

// PaginationOptions holds cursor-based pagination parameters.
type PaginationOptions struct {
	Limit  int
	Cursor string // Encoded cursor (typically base64 of last ID/Time)
}

// ReportRepository defines the storage interface for Report aggregates.
type ReportRepository interface {
	// Get retrieves a report by its unique ID.
	Get(id string) (*model.Report, error)

	// Save persists a report aggregate (Upsert).
	Save(report *model.Report) error

	// List retrieves a list of reports based on filters and pagination.
	List(filters FilterOptions, pagination PaginationOptions) ([]*model.Report, string, error)
	// Returns: reports, nextCursor, error

	// CreateIndexes ensures the necessary database indexes exist.
	CreateIndexes() error
}
