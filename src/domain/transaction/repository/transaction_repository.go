package repository

import "github.com/carddemo/project/src/domain/transaction/model"

// TransactionRepository defines the storage interface for Transaction aggregates.
type TransactionRepository interface {
	Get(id string) (*model.Transaction, error)
	Save(aggregate *model.Transaction) error
	Delete(id string) error
	List() ([]*model.Transaction, error)
}
