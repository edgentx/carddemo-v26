package model

import (
	"time"

	"github.com/carddemo/project/src/domain/shared"
)

// CardPolicy represents the CardPolicy aggregate.
type CardPolicy struct {
	shared.AggregateRoot
	CardType       string    `bson:"card_type" json:"card_type"`
	EffectiveDate  time.Time `bson:"effective_date" json:"effective_date"`
	ExpirationDate time.Time `bson:"expiration_date" json:"expiration_date"`
	Version        int       `bson:"version" json:"version"`
}
