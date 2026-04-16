package model

import (
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// ExportJob represents the aggregate root for export tracking.
type ExportJob struct {
	shared.AggregateBase
	ID        string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	RetryCount int
}

// NewExportJob creates a new ExportJob aggregate.
func NewExportJob(id string) *ExportJob {
	return &ExportJob{
		ID:        id,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		RetryCount: 0,
		AggregateBase: shared.AggregateBase{},
	}
}

// ID method satisfies shared.Aggregate interface via explicit method if needed,
// though typically embedded. We rely on the struct field ID for serialization.
func (e *ExportJob) AggregateID() string {
	return e.ID
}
