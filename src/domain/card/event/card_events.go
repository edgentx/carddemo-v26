package event

import "time"

// CardIssued is emitted when a new card is successfully created.
type CardIssued struct {
	AggregateID string    `json:"aggregate_id"`
	AccountID   string    `json:"account_id"`
	CardType    string    `json:"card_type"`
	IssuedAt    time.Time `json:"issued_at"`
}
