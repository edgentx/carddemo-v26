package event

import (
	"github.com/google/uuid"
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// AccountStatusUpdated is published when an account's status changes.
type AccountStatusUpdated struct {
	shared.CloudEventEnvelope
	Payload struct {
		AccountID string `json:"account_id"`
		OldStatus string `json:"old_status"`
		NewStatus string `json:"new_status"`
		Reason    string `json:"reason"`
	} `json:"data"`
}

// NewAccountStatusUpdated creates a new AccountStatusUpdated event.
func NewAccountStatusUpdated(aggregateID string) *AccountStatusUpdated {
	e := &AccountStatusUpdated{}
	e.ID = uuid.New().String()
	e.Source = "/account"
	e.SpecVersion = "1.0"
	e.Type = "com.carddemo.account.status.updated"
	e.DataContentType = "application/json"
	e.Time = time.Now().Format(time.RFC3339)
	e.Subject = aggregateID
	return e
}

// Type returns the CloudEvent type.
func (e *AccountStatusUpdated) Type() string {
	return e.CloudEventEnvelope.Type
}

// AggregateID returns the aggregate ID.
func (e *AccountStatusUpdated) AggregateID() string {
	return e.Subject
}
