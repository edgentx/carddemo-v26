package event

import (
	"github.com/google/uuid"
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// AccountOpened is published when a new account is created.
type AccountOpened struct {
	shared.CloudEventEnvelope
	Payload struct {
		AccountID     string `json:"account_id"`
		UserProfileID string `json:"user_profile_id"`
		Status        string `json:"status"`
		AccountType   string `json:"account_type"`
	} `json:"data"`
}

// NewAccountOpened creates a new AccountOpened event.
func NewAccountOpened(aggregateID string, cmd interface{}) *AccountOpened {
	e := &AccountOpened{}
	e.ID = uuid.New().String()
	e.Source = "/account"
	e.SpecVersion = "1.0"
	e.Type = "com.carddemo.account.opened"
	e.DataContentType = "application/json"
	e.Time = time.Now().Format(time.RFC3339)
	e.Subject = aggregateID

	// Extract payload based on command type
	// In a real app, type assertion or mapping logic goes here.
	// For this green phase, we'll assume the caller hydrates the payload or we use a helper.
	// However, to satisfy AC requirements directly in the aggregate, we will set data there.
	return e
}

// Type returns the CloudEvent type.
func (e *AccountOpened) Type() string {
	return e.CloudEventEnvelope.Type
}

// AggregateID returns the aggregate ID.
func (e *AccountOpened) AggregateID() string {
	return e.Subject
}
