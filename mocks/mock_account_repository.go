package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
)

// MockAccountRepository is an in-memory implementation of repository.AccountRepository.
type MockAccountRepository struct {
	mu   sync.RWMutex
	data map[string]*model.Account
}

// NewMockAccountRepository creates a new mock repository.
func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{
		data: make(map[string]*model.Account),
	}
}

// Ensure MockAccountRepository implements the interface.
var _ repository.AccountRepository = (*MockAccountRepository)(nil)

// Get retrieves an aggregate by ID.
func (m *MockAccountRepository) Get(id string) (*model.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[id], nil
}

// Save stores an aggregate.
func (m *MockAccountRepository) Save(aggregate *model.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

// Delete removes an aggregate.
func (m *MockAccountRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

// List returns all aggregates.
func (m *MockAccountRepository) List() ([]*model.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	arr := make([]*model.Account, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
