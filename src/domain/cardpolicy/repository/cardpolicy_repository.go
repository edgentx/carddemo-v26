package repository

import "github.com/carddemo/project/src/domain/cardpolicy/model"

// CardPolicyRepository defines the storage interface for CardPolicy aggregates.
type CardPolicyRepository interface {
	Get(id string) (*model.CardPolicy, error)
	Save(aggregate *model.CardPolicy) error
	Delete(id string) error
	List() ([]*model.CardPolicy, error)
}
