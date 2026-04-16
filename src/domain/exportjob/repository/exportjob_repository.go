package repository

import (
	"context"
	"github.com/card-demo/project/src/domain/exportjob/model"
)

// ExportJobRepository defines the persistence interface for Export Jobs.
type ExportJobRepository interface {
	Get(ctx context.Context, id string) (*model.ExportJobAggregate, error)
	Save(ctx context.Context, agg *model.ExportJobAggregate) error
	List(ctx context.Context, filters ListFilters) ([]model.ExportJobAggregate, error)
}