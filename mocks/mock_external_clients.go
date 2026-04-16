package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

// MockTemporalClient wraps the real Temporal client interface or a custom one.
// Based on the feedback, we will use the interface defined in mocks or the real one if dependencies are added.
// Here we define a mock for the relevant methods used by the ExportService.

type MockTemporalClient struct {
	mock.Mock
}

func (m *MockTemporalClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	argList := []interface{}{ctx, options, workflow}
	argList = append(argList, args...)
	args := m.Called(argList...)
	return args.Get(0).(client.WorkflowRun), args.Error(1)
}

// MockStorageClient simulates S3/MinIO interactions
type MockStorageClient struct {
	mock.Mock
}

func (m *MockStorageClient) GetFile(key string) ([]byte, string, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}

func (m *MockStorageClient) PutFile(key string, data []byte) (string, error) {
	args := m.Called(key, data)
	return args.String(0), args.Error(1)
}