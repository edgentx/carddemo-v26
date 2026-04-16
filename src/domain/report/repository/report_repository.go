package repository

import "github.com/carddemo/project/src/domain/report/model"

// ReportRepository defines the storage interface for Report aggregates.
type ReportRepository interface {
	Get(id string) (*model.Report, error)
	Save(aggregate *model.Report) error
	Delete(id string) error
	List() ([]*model.Report, error)
}
