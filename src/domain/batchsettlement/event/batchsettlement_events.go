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
