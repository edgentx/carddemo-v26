package service

import (
	"context"
	"errors"

	"github.com/carddemo/project/src/domain/cardpolicy/model"
	"github.com/carddemo/project/src/domain/cardpolicy/repository"
)

var (
	// ErrPolicyNotFound is returned when a policy cannot be found.
	ErrPolicyNotFound = errors.New("policy not found")
)

// CardPolicyApplicationService handles policy use cases.
type CardPolicyApplicationService struct {
	policyRepo repository.CardPolicyRepository
}

// NewCardPolicyApplicationService creates a new CardPolicyApplicationService.
func NewCardPolicyApplicationService(policyRepo repository.CardPolicyRepository) *CardPolicyApplicationService {
	return &CardPolicyApplicationService{
		policyRepo: policyRepo,
	}
}

// GetPolicy retrieves a policy by ID.
func (s *CardPolicyApplicationService) GetPolicy(ctx context.Context, id string) (*model.CardPolicy, error) {
	policy, err := s.policyRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, ErrPolicyNotFound
	}
	return policy, nil
}

// UpdatePolicyLimits updates the daily limit of a policy.
func (s *CardPolicyApplicationService) UpdatePolicyLimits(ctx context.Context, id string, dailyLimit int) error {
	policy, err := s.policyRepo.Get(id)
	if err != nil {
		return err
	}
	if policy == nil {
		return ErrPolicyNotFound
	}

	policy.DailyLimit = dailyLimit
	return s.policyRepo.Save(policy)
}

// ListPoliciesByAccount retrieves all policies for an account.
func (s *CardPolicyApplicationService) ListPoliciesByAccount(ctx context.Context, accountID string) ([]*model.CardPolicy, error) {
	// Since the repository interface only has List() (no query by account),
	// we must fetch all and filter in memory for this prototype.
	// In a real Mongo repo, this would be a specific query.
	allPolicies, err := s.policyRepo.List()
	if err != nil {
		return nil, err
	}

	var result []*model.CardPolicy
	for _, p := range allPolicies {
		if p.AccountID == accountID {
			result = append(result, p)
		}
	}
	return result, nil
}
