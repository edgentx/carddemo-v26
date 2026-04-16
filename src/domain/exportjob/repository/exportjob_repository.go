package repository

import (
	"context"

	"github.com/carddemo/project/src/domain/exportjob/model"
)

// ExportJobRepository defines the storage interface for ExportJob aggregates.
type ExportJobRepository interface {
	// Get retrieves an export job by its unique ID.
	Get(id string) (*model.ExportJob, error)

	// Save persists an export job aggregate.
	Save(job *model.ExportJob) error

	// UpdateStatus atomically updates the status of an export job.
	UpdateStatus(id string, newStatus string) error

	// IncrementRetry atomically increments the retry count for an export job.
	IncrementRetry(id string) error

	// CreateIndexes ensures the necessary database indexes exist.
	CreateIndexes() error
}
