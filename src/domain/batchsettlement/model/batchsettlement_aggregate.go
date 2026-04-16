package model

import (
	"time"

	"github.com/google/uuid"
)

// BatchSettlement represents the domain aggregate for grouping transactions for settlement.
type BatchSettlement struct {
	ID          string
	MerchantID  string
	Date        time.Time
	Status      string // e.g., "open", "reconciled"
	TransactionIds []string
	CreatedAt   time.Time
	Version     int
}

// SettlementGroup represents the result of a settlement aggregation query.
type SettlementGroup struct {
	MerchantID string
	Date       time.Time
	Count      int
	Total      float64
}

// NewBatchSettlement creates a new BatchSettlement aggregate.
func NewBatchSettlement(merchantID string, date time.Time) *BatchSettlement {
	return &BatchSettlement{
		ID:          uuid.New().String(),
		MerchantID:  merchantID,
		Date:        date,
		Status:      "open",
		TransactionIds: []string{},
		CreatedAt:   time.Now(),
		Version:     0,
	}
}
