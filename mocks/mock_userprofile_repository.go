package mocks

import (
	"sync"

	"github.com/carddemo/project/src/domain/userprofile/model"
	profileRepo "github.com/carddemo/project/src/domain/userprofile/repository"
)

// MockUserProfileRepository is an in-memory implementation for testing.
type MockUserProfileRepository struct {
	mu   sync.RWMutex
	data map[string]*model.UserProfile
}

func NewMockUserProfileRepository() *MockUserProfileRepository {
	return &MockUserProfileRepository{data: make(map[string]*model.UserProfile)}
}

var _ profileRepo.UserProfileRepository = (*MockUserProfileRepository)(nil)

func (m *MockUserProfileRepository) Get(id string) (*model.UserProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if val, ok := m.data[id]; ok {
		return val, nil
	}
	return nil, nil // Return nil, nil to simulate not found
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
