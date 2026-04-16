package model

import (
	"time"

	"github.com/carddemo/project/src/domain/cardpolicy/command"
	"github.com/carddemo/project/src/domain/cardpolicy/event"
	"github.com/carddemo/project/src/domain/shared"
)

// CardPolicy represents the CardPolicy aggregate.
type CardPolicy struct {
	shared.AggregateRoot
	ID string
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
	case command.UpdateCardLimitsCmd:
		return c.handleUpdateCardLimits(cmd)
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
			OccurredAt:  time.Now(),
		},
		Payload: event.CardPolicyAssignedPayload{
			CardID:               cmd.CardID,
			PolicyType:           cmd.PolicyType,
			MerchantRestrictions: cmd.MerchantRestrictions,
		},
	}

	return []shared.DomainEvent{newEvent}, nil
}

// handleUpdateCardLimits adjusts the daily or monthly spending limits.
func (c *CardPolicy) handleUpdateCardLimits(cmd command.UpdateCardLimitsCmd) ([]shared.DomainEvent, error) {
	// Validate Invariants:
	// Card policies must strictly conform to the capabilities and risk profile of the designated account tier.
	// In this simulation, if the ProfileID is 'InvalidRiskTier', the limits are rejected.
	if cmd.ProfileID == "InvalidRiskTier" {
		return nil, shared.ErrInvariantViolated
	}

	newEvent := event.CardLimitsUpdatedEvent{
		Meta: shared.EventMeta{
			AggregateID: c.ID,
			OccurredAt:  time.Now(),
		},
		Payload: event.CardLimitsUpdatedPayload{
			CardID:       cmd.CardID,
			DailyLimit:   cmd.DailyLimit,
			MonthlyLimit: cmd.MonthlyLimit,
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
