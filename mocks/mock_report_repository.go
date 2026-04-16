package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/report/model"
	"github.com/carddemo/project/src/domain/report/repository"
)

type MockReportRepository struct {
	mu   sync.RWMutex
	data map[string]*model.Report
}

func NewMockReportRepository() *MockReportRepository {
	return &MockReportRepository{data: make(map[string]*model.Report)}
}

var _ repository.ReportRepository = (*MockReportRepository)(nil)

func (m *MockReportRepository) Get(id string) (*model.Report, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[id], nil
}

func (m *MockReportRepository) Save(aggregate *model.Report) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

func (m *MockReportRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

func (m *MockReportRepository) List() ([]*model.Report, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]*model.Report, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
