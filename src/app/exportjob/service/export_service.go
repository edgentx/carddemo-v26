package service

import (
	"context"
	"fmt"

	"github.com/card-demo/project/src/domain/exportjob/command"
	"github.com/card-demo/project/src/domain/exportjob/model"
	"github.com/card-demo/project/src/domain/exportjob/repository"
)

// TemporalClient defines the interface for interacting with Temporal.
// This is a local port interface to avoid direct dependency on the concrete SDK in the app layer.
type TemporalClient interface {
	ExecuteWorkflow(ctx context.Context, workflow interface{}, args ...interface{}) (string, error)
}

// ExportService handles the use cases for Export Jobs.
type ExportService struct {
	repo           repository.ExportJobRepository
	temporalClient TemporalClient
}

// NewExportService creates a new ExportService.
func NewExportService(repo repository.ExportJobRepository, temporalClient TemporalClient) *ExportService {
	return &ExportService{
		repo:           repo,
		temporalClient: temporalClient,
	}
}

// CreateExport initiates a new export job and triggers the async workflow.
func (s *ExportService) CreateExport(ctx context.Context, reportID string) (*model.ExportJobAggregate, error) {
	// 1. Create the Aggregate
	agg := model.NewExportJobAggregate(reportID)

	// 2. Handle the command to transition state
	err := agg.Handle(command.InitiateExport{
		ReportID: reportID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initiate export: %w", err)
	}

	// 3. Persist initial state
	err = s.repo.Save(ctx, agg)
	if err != nil {
		return nil, fmt.Errorf("failed to save export job: %w", err)
	}

	// 4. Trigger Temporal Workflow (Fire and Forget)
	// We pass the aggregate ID so the workflow can update the status later.
	_, err = s.temporalClient.ExecuteWorkflow(ctx, "export-workflow", agg.ID)
	if err != nil {
		// In a real system, we might want to mark the aggregate as failed or schedule a retry.
		// For now, we log the error but the job is created.
		return nil, fmt.Errorf("failed to trigger workflow: %w", err)
	}

	return agg, nil
}

// Get retrieves an export job by ID.
func (s *ExportService) Get(ctx context.Context, id string) (*model.ExportJobAggregate, error) {
	return s.repo.Get(ctx, id)
}
