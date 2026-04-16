package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/exportjob/model"
	"github.com/carddemo/project/src/domain/exportjob/repository"
)

type MockExportJobRepository struct {
	mu   sync.RWMutex
	data map[string]*model.ExportJob
}

func NewMockExportJobRepository() *MockExportJobRepository {
	return &MockExportJobRepository{data: make(map[string]*model.ExportJob)}
}

var _ repository.ExportJobRepository = (*MockExportJobRepository)(nil)

func (m *MockExportJobRepository) Get(id string) (*model.ExportJob, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[id], nil
}

func (m *MockExportJobRepository) Save(aggregate *model.ExportJob) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

func (m *MockExportJobRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

func (m *MockExportJobRepository) List() ([]*model.ExportJob, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]*model.ExportJob, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
