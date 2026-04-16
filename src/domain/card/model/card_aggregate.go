package model

import (
	"github.com/carddemo/project/src/domain/shared"
)

// Card represents the Card aggregate.
type Card struct {
	shared.AggregateRoot
	AccountID        string `bson:"account_id" json:"account_id"`
	CardNumber       string `bson:"card_number" json:"-"` // Never return full PAN in JSON
	HashedCardNumber string `bson:"hashed_card_number" json:"-"`
	Status           string `bson:"status" json:"status"`
	Type             string `bson:"type" json:"type"`
	Version          int    `bson:"version" json:"version"`
}
