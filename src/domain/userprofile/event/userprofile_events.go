package event

import (
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// UserRegistered is emitted when a new user profile is created.
type UserRegistered struct {
	shared.DomainEventBase
	AggregateID string `json:"aggregate_id"`
	Email       string `json:"email"`
	CreditScore int    `json:"credit_score"`
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

// UserLinkedToAccount is emitted when a user is linked to an account.
type UserLinkedToAccount struct {
	shared.DomainEventBase
	AggregateID string `json:"aggregate_id"`
	AccountID   string `json:"account_id"`
}

// NewUserLinkedToAccount creates a UserLinkedToAccount event.
func NewUserLinkedToAccount(aggregateID string, accountID string) shared.DomainEvent {
	return &UserLinkedToAccount{
		DomainEventBase: shared.DomainEventBase{
			Type:      "com.carddemo.user.linked.to.account",
			Timestamp: time.Now().UTC(),
		},
		AggregateID: aggregateID,
		AccountID:   accountID,
	}
}