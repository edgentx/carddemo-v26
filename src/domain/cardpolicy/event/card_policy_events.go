package event

import (
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// CardPolicyAssigned is emitted when a policy is successfully assigned to a card.
type CardPolicyAssigned struct {
	Meta    shared.EventMeta
	Payload CardPolicyAssignedPayload
}

// Type returns the CloudEvents type for this event.
func (e CardPolicyAssigned) Type() string {
	return "com.carddemo.cardpolicy.assigned"
}

// OccurredAt returns when the event happened.
func (e CardPolicyAssigned) OccurredAt() time.Time {
	return e.Meta.OccurredAt
}

// AggregateID returns the ID of the aggregate that emitted this event.
func (e CardPolicyAssigned) AggregateID() string {
	return e.Meta.AggregateID
}

// GetPayload returns the payload of the event.
func (e CardPolicyAssigned) GetPayload() interface{} {
	return e.Payload
}

type CardPolicyAssignedPayload struct {
	CardID               string   `json:"card_id"`
	PolicyType           string   `json:"policy_type"`
	MerchantRestrictions []string `json:"merchant_restrictions"`
}
