package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/cardpolicy/model"
	"github.com/carddemo/project/src/domain/cardpolicy/repository"
)

type MockCardPolicyRepository struct {
	mu   sync.RWMutex
	data map[string]*model.CardPolicy
}

func NewMockCardPolicyRepository() *MockCardPolicyRepository {
	return &MockCardPolicyRepository{data: make(map[string]*model.CardPolicy)}
}

var _ repository.CardPolicyRepository = (*MockCardPolicyRepository)(nil)

func (m *MockCardPolicyRepository) Get(id string) (*model.CardPolicy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[id], nil
}

func (m *MockCardPolicyRepository) Save(aggregate *model.CardPolicy) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

func (m *MockCardPolicyRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

func (m *MockCardPolicyRepository) List() ([]*model.CardPolicy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	arr := make([]*model.CardPolicy, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
