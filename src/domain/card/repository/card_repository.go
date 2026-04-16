package repository

import "github.com/carddemo/project/src/domain/card/model"

// CardRepository defines the storage interface for Card aggregates.
type CardRepository interface {
	Get(id string) (*model.Card, error)
	Save(aggregate *model.Card) error
	Delete(id string) error
	List() ([]*model.Card, error)
}
