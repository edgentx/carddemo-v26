package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
)

// MockAccountRepository is an in-memory implementation for testing.
type MockAccountRepository struct {
	mu   sync.RWMutex
	data map[string]*model.Account
}

func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{data: make(map[string]*model.Account)}
}

var _ repository.AccountRepository = (*MockAccountRepository)(nil)

func (m *MockAccountRepository) Get(id string) (*model.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if val, ok := m.data[id]; ok {
		return val, nil
	}
	return nil, nil // Return nil, nil to simulate not found or specific error if desired
}

func (m *MockAccountRepository) Save(aggregate *model.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

func (m *MockAccountRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

func (m *MockAccountRepository) List() ([]*model.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]*model.Account, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
