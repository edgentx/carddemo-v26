package event

import (
	"github.com/carddemo/project/src/domain/shared"
	"time"
)

// UserRegistered is emitted when a new user profile is created.
type UserRegistered struct {
	shared.DomainEventBase
	AggregateID string `json:"aggregate_id"`
	Email string `json:"email"`
	CreditScore int `json:"credit_score"`
}

// NewUserRegistered creates a UserRegistered event.
func NewUserRegistered(aggregateID string, email string, creditScore int) shared.DomainEvent {
	return &UserRegistered{
		DomainEventBase: shared.DomainEventBase{
			Type:      "com.carddemo.user.registered",
			Timestamp: time.Now().UTC(),
			// Data will be populated by the aggregate
		},
		AggregateID: aggregateID,
		Email:       email,
		CreditScore: creditScore,
	}
}
