package model_test

import (
	"testing"

	"github.com/carddemo/project/src/domain/cardpolicy/command"
	"github.com/carddemo/project/src/domain/cardpolicy/model"
	"github.com/carddemo/project/src/domain/shared"
)

func TestCardPolicy_AssignCardPolicyCmd(t *testing.T) {
	tests := []struct {
		name        string
		aggregate   *model.CardPolicy
		cmd         command.AssignCardPolicyCmd
		wantEvents  int
		wantErr     error
		eventType   string
		policyType  string
		cardID      string
	}{
		{
			name:      "Successfully execute AssignCardPolicyCmd",
			aggregate: model.NewCardPolicy("cp-123"),
			cmd: command.AssignCardPolicyCmd{
				CardID:               "card-456",
				PolicyType:           "Standard",
				MerchantRestrictions: []string{"Travel", "Dining"},
			},
			wantEvents: 1,
			wantErr:    nil,
			eventType:  "com.carddemo.cardpolicy.assigned",
			policyType: "Standard",
			cardID:     "card-456",
		},
		{
			name:      "AssignCardPolicyCmd rejected - Invalid Risk Tier",
			aggregate: model.NewCardPolicy("cp-789"),
			cmd: command.AssignCardPolicyCmd{
				CardID:               "card-999",
				PolicyType:           "InvalidRiskTier",
				MerchantRestrictions: []string{},
			},
			wantEvents: 0,
			wantErr:    shared.ErrInvariantViolated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := tt.aggregate.Execute(tt.cmd)

			if err != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(events) != tt.wantEvents {
				t.Errorf("Execute() events count = %v, wantEvents %v", len(events), tt.wantEvents)
				return
			}

			if tt.wantEvents > 0 {
				evt, ok := events[0].(interface{ Type() string })
				if !ok {
					t.Error("Event does not implement Type() method")
					return
				}

				if evt.Type() != tt.eventType {
					t.Errorf("Event type = %v, want %v", evt.Type(), tt.eventType)
				}

				// Verify payload fields match
				payloadEvents, ok := events[0].(interface{ GetPayload() interface{} })
				if !ok {
					t.Error("Event does not implement GetPayload() method")
					return
				}
				payload := payloadEvents.GetPayload()
				p, ok := payload.(struct {
					CardID               string   `json:"card_id"`
					PolicyType           string   `json:"policy_type"`
					MerchantRestrictions []string `json:"merchant_restrictions"`
				})
				if !ok {
					t.Errorf("Payload structure mismatch. Got %+v", payload)
					return
				}

				if p.CardID != tt.cardID {
					t.Errorf("Payload CardID = %v, want %v", p.CardID, tt.cardID)
				}
				if p.PolicyType != tt.policyType {
					t.Errorf("Payload PolicyType = %v, want %v", p.PolicyType, tt.policyType)
				}
			}
		})
	}
}
