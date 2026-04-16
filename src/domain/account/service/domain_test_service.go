package service

import (
	"errors"
	"github.com/carddemo/project/src/domain/account/command"
	"github.com/carddemo/project/src/domain/account/model"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/shared"
	"github.com/google/uuid"
)

// AccountDomainService is a lightweight orchestrator if needed,
// but in this architecture, the Aggregates are king.
// This service acts as the entry point for testing the Red Phase behavior
// by wiring the Repository to the Aggregate.

type AccountDomainService struct {
	Repo repository.AccountRepository
}

func NewAccountDomainService(repo repository.AccountRepository) *AccountDomainService {
	return &AccountDomainService{Repo: repo}
}

// ExecuteCommand simulates the UseCase layer behavior for testing purposes.
func (s *AccountDomainService) ExecuteCommand(aggregateID string, cmd interface{}) (*model.AccountAggregate, error) {
	var agg *model.AccountAggregate
	var err error

	if aggregateID != "" {
		agg, err = s.Repo.Get(aggregateID)
		if err != nil {
			return nil, err
		}
	}

	// Handle specific commands
	switch c := cmd.(type) {
	case *command.OpenAccountCmd:
		// New Aggregate
		newID := uuid.New().String()
		agg = model.NewAccountAggregate(newID)
		err = agg.Handle(c)
	case *command.UpdateAccountStatusCmd:
		// Existing Aggregate
		err = agg.Handle(c)
	default:
		return nil, errors.New("unknown command")
	}

	if err != nil {
		return nil, err
	}

	// Persist
	saveErr := s.Repo.Save(agg)
	if saveErr != nil {
		return nil, saveErr
	}

	return agg, nil
}

// --- Interface & Mock Definitions for Tests ---

// MockAccountRepository is a memory-based implementation of AccountRepository for testing.
type MockAccountRepository struct {
	Data map[string]*model.AccountAggregate
}

func NewMockAccountRepository() *MockAccountRepository {
	return &MockAccountRepository{
		Data: make(map[string]*model.AccountAggregate),
	}
}

func (m *MockAccountRepository) Get(id string) (*model.AccountAggregate, error) {
	if agg, ok := m.Data[id]; ok {
		// Return a copy to prevent direct state mutation outside of transaction
		// Deep copy logic omitted for brevity, but conceptually important.
		return agg, nil
	}
	return nil, shared.ErrAggregateNotFound
}

func (m *MockAccountRepository) Save(agg *model.AccountAggregate) error {
	// Optimistic Lock Check would go here in real impl
	m.Data[agg.ID] = agg
	return nil
}

func (m *MockAccountRepository) Delete(id string) error {
	delete(m.Data, id)
	return nil
}

func (m *MockAccountRepository) List() ([]*model.AccountAggregate, error) {
	var list []*model.AccountAggregate
	for _, v := range m.Data {
		list = append(list, v)
	}
	return list, nil
}

var _ repository.AccountRepository = (*MockAccountRepository)(nil)
