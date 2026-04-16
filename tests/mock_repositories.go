package mocks

import (
	"context"

	"github.com/card-demo/project/src/domain/report/repository"
	"github.com/card-demo/project/src/domain/exportjob/repository"
	"github.com/card-demo/project/src/domain/report/model"
	export_model "github.com/card-demo/project/src/domain/exportjob/model"
	"github.com/stretchr/testify/mock"
)

// MockReportRepository
type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) Get(ctx context.Context, id string) (*model.ReportAggregate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReportAggregate), args.Error(1)
}

func (m *MockReportRepository) Save(ctx context.Context, agg *model.ReportAggregate) error {
	args := m.Called(ctx, agg)
	return args.Error(0)
}

func (m *MockReportRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockReportRepository) List(ctx context.Context, filters repository.ListFilters) ([]model.ReportAggregate, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ReportAggregate), args.Error(1)
}

// MockExportJobRepository
type MockExportJobRepository struct {
	mock.Mock
}

func (m *MockExportJobRepository) Get(ctx context.Context, id string) (*export_model.ExportJobAggregate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*export_model.ExportJobAggregate), args.Error(1)
}

func (m *MockExportJobRepository) Save(ctx context.Context, agg *export_model.ExportJobAggregate) error {
	args := m.Called(ctx, agg)
	return args.Error(0)
}

func (m *MockExportJobRepository) List(ctx context.Context, filters repository.ListFilters) ([]export_model.ExportJobAggregate, error) {
	// Assuming similar list filters for now
	args := m.Called(ctx, filters)
	return args.Get(0).([]export_model.ExportJobAggregate), args.Error(1)
}