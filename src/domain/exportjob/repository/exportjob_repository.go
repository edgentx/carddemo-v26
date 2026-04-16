package repository

import "github.com/carddemo/project/src/domain/exportjob/model"

// ExportJobRepository defines the storage interface for ExportJob aggregates.
type ExportJobRepository interface {
	Get(id string) (*model.ExportJob, error)
	Save(aggregate *model.ExportJob) error
	Delete(id string) error
	List() ([]*model.ExportJob, error)
}
