package model

import (
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// Report represents the aggregate root for financial reports.
type Report struct {
	shared.AggregateBase
	ID        string
	AccountID string
	Type      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      map[string]interface{}
}

// NewReport creates a new Report aggregate.
func NewReport(id, accountID, reportType string) *Report {
	return &Report{
		ID:        id,
		AccountID: accountID,
		Type:      reportType,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		AggregateBase: shared.AggregateBase{},
	}
}

// AggregateID satisfies shared.Aggregate interface
func (r *Report) AggregateID() string {
	return r.ID
}
