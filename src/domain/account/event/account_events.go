package event

import (
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
