package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/src/domain/card/repository"
)

type MockCardRepository struct {
	mu   sync.RWMutex
	data map[string]*model.Card
}

func NewMockCardRepository() *MockCardRepository {
	return &MockCardRepository{data: make(map[string]*model.Card)}
}

var _ repository.CardRepository = (*MockCardRepository)(nil)

func (m *MockCardRepository) Get(id string) (*model.Card, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[id], nil
}

func (m *MockCardRepository) Save(aggregate *model.Card) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

func (m *MockCardRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

func (m *MockCardRepository) List() ([]*model.Card, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]*model.Card, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
