package repository

import (
	"github.com/carddemo/project/src/domain/batchsettlement/model"
)

// BatchSettlementRepository defines the storage interface for Batch Settlements.
type BatchSettlementRepository interface {
	Get(id string) (*model.BatchSettlement, error)
	Save(aggregate *model.BatchSettlement) error
	List() ([]*model.BatchSettlement, error)
}
