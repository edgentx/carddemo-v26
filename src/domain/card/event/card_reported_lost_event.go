package event

import "time"

// CardReportedLost is emitted when a card is successfully reported as lost.
type CardReportedLost struct {
	AggregateID string    `json:"aggregate_id"`
	Reason      string    `json:"reason"`
	ReportedBy  string    `json:"reported_by"`
	OccurredAt  time.Time `json:"occurred_at"`
}
