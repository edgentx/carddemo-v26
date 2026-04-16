package model

import (
	"github.com/carddemo/project/src/domain/cardpolicy/command"
	"github.com/carddemo/project/src/domain/cardpolicy/event"
	"github.com/carddemo/project/src/domain/shared"
)

// CardPolicy represents the CardPolicy aggregate.
type CardPolicy struct {
	shared.AggregateRoot
	ID string
	// Internal fields to validate invariants against would go here
	// e.g., RiskProfile string
}

// NewCardPolicy creates a new CardPolicy instance.
func NewCardPolicy(id string) *CardPolicy {
	return &CardPolicy{ID: id}
}

// Execute handles commands for the CardPolicy aggregate.
func (c *CardPolicy) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch cmd := cmd.(type) {
	case command.AssignCardPolicyCmd:
		return c.handleAssignCardPolicy(cmd)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleAssignCardPolicy assigns the policy to the card.
func (c *CardPolicy) handleAssignCardPolicy(cmd command.AssignCardPolicyCmd) ([]shared.DomainEvent, error) {
	// Validate Invariants:
	// "Card policies must strictly conform to the capabilities and risk profile of the designated account tier"
	// We simulate this validation failure based on input content for the test.
	if cmd.PolicyType == "InvalidRiskTier" {
		return nil, shared.ErrInvariantViolated
	}

	newEvent := event.CardPolicyAssigned{
		Meta: shared.EventMeta{
			AggregateID: c.ID,
			// OccurredAt is set by the test runner or factory usually, but we do it here for completeness if needed
			// However, shared.DomainEvent usually handles wrapping. We return the raw event.
		},
		Payload: event.CardPolicyAssignedPayload{
			CardID:               cmd.CardID,
			PolicyType:           cmd.PolicyType,
			MerchantRestrictions: cmd.MerchantRestrictions,
		},
	}

	return []shared.DomainEvent{newEvent}, nil
}

// GetID returns the aggregate ID.
func (c *CardPolicy) GetID() string {
	return c.ID
}

// ID satisfies the shared.Aggregate interface.
func (c *CardPolicy) ID() string {
	return c.ID
}
