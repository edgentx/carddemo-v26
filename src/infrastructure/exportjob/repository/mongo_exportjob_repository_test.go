package repository

import (
	"context"
	"testing"
	"time"

	"github.com/carddemo/project/src/domain/exportjob/model"
	"github.com/carddemo/project/src/infrastructure/shared/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMongoExportJobRepository_SaveAndGet(t *testing.T) {
	mockClient := mocks.NewMockMongoClient()
	repo := NewMongoExportJobRepository(mockClient)

	job := &model.ExportJob{
		ID:        "job_123",
		Status:    "pending",
		CreatedAt: time.Now(),
		RetryCount: 0,
	}

	err := repo.Save(job)
	require.NoError(t, err)

	retrieved, err := repo.Get("job_123")
	require.NoError(t, err)
	assert.Equal(t, job.ID, retrieved.ID)
	assert.Equal(t, job.Status, retrieved.Status)
}

func TestMongoExportJobRepository_UpdateStatus(t *testing.T) {
	mockClient := mocks.NewMockMongoClient()
	repo := NewMongoExportJobRepository(mockClient)

	// Setup: Save initial state
	job := &model.ExportJob{ID: "job_status", Status: "pending", RetryCount: 0}
	_ = repo.Save(job)

	// Execute: Update Status
	err := repo.UpdateStatus("job_status", "processing")
	require.NoError(t, err, "UpdateStatus should succeed")

	// Verify: Fetch and check status
	updated, err := repo.Get("job_status")
	require.NoError(t, err)
	assert.Equal(t, "processing", updated.Status)
}

func TestMongoExportJobRepository_IncrementRetry(t *testing.T) {
	mockClient := mocks.NewMockMongoClient()
	repo := NewMongoExportJobRepository(mockClient)

	job := &model.ExportJob{ID: "job_retry", Status: "failed", RetryCount: 1}
	_ = repo.Save(job)

	err := repo.IncrementRetry("job_retry")
	require.NoError(t, err, "IncrementRetry should succeed")

	updated, err := repo.Get("job_retry")
	require.NoError(t, err)
	assert.Equal(t, 2, updated.RetryCount, "Retry count should be incremented by 1")
}

func TestMongoExportJobRepository_CreateIndexes(t *testing.T) {
	ctx := context.Background()
	mockClient := mocks.NewMockMongoClient()
	repo := NewMongoExportJobRepository(mockClient)

	err := repo.CreateIndexes(ctx)
	assert.NoError(t, err)
}
