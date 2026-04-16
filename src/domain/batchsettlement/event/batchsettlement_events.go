package event

import "github.com/carddemo/project/src/domain/shared"

// BatchOpened is emitted when a new settlement batch is successfully initiated.
type BatchOpened struct {
	shared.DomainEventMeta
	SettlementDate    string `json:"settlement_date"`
	OperationalRegion string `json:"operational_region"`
}

// Type implements DomainEvent.
func (e BatchOpened) Type() string {
	return "com.carddemo.batch.opened"
}

// BatchReconciled is emitted when a batch is successfully reconciled and frozen.
type BatchReconciled struct {
	shared.DomainEventMeta
	TotalDebits  int64 `json:"total_debits"`
	TotalCredits int64 `json:"total_credits"`
}

// Type implements DomainEvent.
func (e BatchReconciled) Type() string {
	return "com.carddemo.batch.reconciled"
}
