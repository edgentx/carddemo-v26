package service

import (
	"context"
	"fmt"

	"github.com/card-demo/project/src/domain/report/command"
	"github.com/card-demo/project/src/domain/report/model"
	"github.com/card-demo/project/src/domain/report/repository"
)

// ReportService handles the use cases for Reports.
type ReportService struct {
	repo repository.ReportRepository
}

// NewReportService creates a new ReportService.
func NewReportService(repo repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// CreateReport creates a new report request.
func (s *ReportService) CreateReport(ctx context.Context, reportType string, params map[string]string) (*model.ReportAggregate, error) {
	// 1. Create Aggregate
	agg := model.NewReportAggregate(reportType, params)

	// 2. Handle Command
	err := agg.Handle(command.RequestReport{
		Type:   reportType,
		Params: params,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to request report: %w", err)
	}

	// 3. Save
	err = s.repo.Save(ctx, agg)
	if err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	return agg, nil
}

// Get retrieves a report by ID.
func (s *ReportService) Get(ctx context.Context, id string) (*model.ReportAggregate, error) {
	return s.repo.Get(ctx, id)
}

// List retrieves reports based on filters.
func (s *ReportService) List(ctx context.Context, filters repository.ListFilters) ([]model.ReportAggregate, error) {
	return s.repo.List(ctx, filters)
}

// Delete removes a report by ID.
func (s *ReportService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
