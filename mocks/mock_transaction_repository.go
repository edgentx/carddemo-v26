package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/transaction/model"
	"github.com/carddemo/project/src/domain/transaction/repository"
)

type MockTransactionRepository struct {
	mu   sync.RWMutex
	data map[string]*model.Transaction
}

func NewMockTransactionRepository() *MockTransactionRepository {
	return &MockTransactionRepository{data: make(map[string]*model.Transaction)}
}

var _ repository.TransactionRepository = (*MockTransactionRepository)(nil)

func (m *MockTransactionRepository) Get(id string) (*model.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[id], nil
}

func (m *MockTransactionRepository) Save(aggregate *model.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

func (m *MockTransactionRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

func (m *MockTransactionRepository) List() ([]*model.Transaction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]*model.Transaction, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
