package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
)

// MockUserProfileRepository is an in-memory implementation of repository.UserProfileRepository.
type MockUserProfileRepository struct {
	mu   sync.RWMutex
	data map[string]*model.UserProfile
}

// NewMockUserProfileRepository creates a new mock repository.
func NewMockUserProfileRepository() *MockUserProfileRepository {
	return &MockUserProfileRepository{
		data: make(map[string]*model.UserProfile),
	}
}

// Ensure MockUserProfileRepository implements the interface.
var _ repository.UserProfileRepository = (*MockUserProfileRepository)(nil)

func (m *MockUserProfileRepository) Get(id string) (*model.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[id], nil
}

func (m *MockUserProfileRepository) Save(aggregate *model.UserProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[aggregate.ID] = aggregate
	return nil
}

func (m *MockUserProfileRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

func (m *MockUserProfileRepository) List() ([]*model.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	arr := make([]*model.UserProfile, 0, len(m.data))
	for _, v := range m.data {
		arr = append(arr, v)
	}
	return arr, nil
}
