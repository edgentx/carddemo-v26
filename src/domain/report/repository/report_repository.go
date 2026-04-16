package repository

import (
	"context"
	"github.com/card-demo/project/src/domain/report/model"
)

type ListFilters struct {
	Type      string
	Status    string
	StartDate string
	EndDate   string
	Page      int
	Limit     int
}

// ReportRepository defines the persistence interface for Reports.
type ReportRepository interface {
	Get(ctx context.Context, id string) (*model.ReportAggregate, error)
	Save(ctx context.Context, agg *model.ReportAggregate) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters ListFilters) ([]model.ReportAggregate, error)
}