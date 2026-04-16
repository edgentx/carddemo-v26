package event

import (
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// CardLimitsUpdatedEvent is emitted when card limits are successfully changed.
type CardLimitsUpdatedEvent struct {
	Meta    shared.EventMeta
	Payload CardLimitsUpdatedPayload
}

// Type returns the CloudEvents type for this event.
func (e CardLimitsUpdatedEvent) Type() string {
	return "com.carddemo.cardpolicy.limits.updated"
}

// OccurredAt returns when the event happened.
func (e CardLimitsUpdatedEvent) OccurredAt() time.Time {
	return e.Meta.OccurredAt
}

// AggregateID returns the ID of the aggregate that emitted this event.
func (e CardLimitsUpdatedEvent) AggregateID() string {
	return e.Meta.AggregateID
}

// GetPayload returns the payload of the event.
func (e CardLimitsUpdatedEvent) GetPayload() interface{} {
	return e.Payload
}

type CardLimitsUpdatedPayload struct {
	CardID       string `json:"card_id"`
	DailyLimit   int64  `json:"daily_limit"`
	MonthlyLimit int64  `json:"monthly_limit"`
}
