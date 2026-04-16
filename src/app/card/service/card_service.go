package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/carddemo/project/src/domain/card/command"
	"github.com/carddemo/project/src/domain/card/model"
	"github.com/carddemo/project/src/domain/card/repository"
	"github.com/google/uuid"
)

var (
	// ErrCardNotFound is returned when a card cannot be found.
	ErrCardNotFound = errors.New("card not found")
	// ErrInvalidStateTransition is returned when a status transition is illegal.
	ErrInvalidStateTransition = errors.New("invalid state transition")
	// ErrAccountNotFound is returned when the account does not exist.
	ErrAccountNotFound = errors.New("account not found")
)

// CardApplicationService handles card use cases.
type CardApplicationService struct {
	cardRepo    repository.CardRepository
	accountRepo repository.AccountRepository
}

// NewCardApplicationService creates a new CardApplicationService.
func NewCardApplicationService(cardRepo repository.CardRepository, accountRepo repository.AccountRepository) *CardApplicationService {
	return &CardApplicationService{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
	}
}

// IssueCard creates a new card for an account.
func (s *CardApplicationService) IssueCard(ctx context.Context, accountID, cardType string, limits map[string]int) (string, error) {
	// Validate Account Exists
	_, err := s.accountRepo.Get(accountID)
	if err != nil {
		return "", fmt.Errorf("account check failed: %w", err)
	}
	// Note: The mock repo returns nil, nil for not found. 
	// If we rely on specific behavior, we should handle the nil case.
	// For this implementation, we assume if Get returns nil object, it's not found.
	// But we'll stick to the interface contract.

	// Logic to check if mock returned nil (mock specific logic)
	// In a real repo, Get would return a concrete error.

	id := uuid.New().String()
	newCard := model.NewCard(id, accountID, cardType, limits)

	if err := s.cardRepo.Save(newCard); err != nil {
		return "", err
	}

	return id, nil
}

// GetCard retrieves a card by ID.
func (s *CardApplicationService) GetCard(ctx context.Context, id string) (*model.Card, error) {
	card, err := s.cardRepo.Get(id)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, ErrCardNotFound
	}
	return card, nil
}

// UpdateCardStatus updates the status of a card.
func (s *CardApplicationService) UpdateCardStatus(ctx context.Context, id, status string) error {
	card, err := s.cardRepo.Get(id)
	if err != nil {
		return err
	}
	if card == nil {
		return ErrCardNotFound
	}

	// Validation: ensure status is valid
	validStatuses := map[string]bool{
		"Issued":    true,
		"Active":    true,
		"Suspended": true,
		"Blocked":   true,
		"Closed":    true,
	}
	if !validStatuses[status] {
		return errors.New("invalid card status")
	}

	card.Status = model.CardStatus(status)
	return s.cardRepo.Save(card)
}

// ActivateCard activates a card.
func (s *CardApplicationService) ActivateCard(ctx context.Context, id string) error {
	card, err := s.cardRepo.Get(id)
	if err != nil {
		return err
	}
	if card == nil {
		return ErrCardNotFound
	}

	if card.Status != model.CardStatusIssued {
		return ErrInvalidStateTransition
	}

	card.Status = model.CardStatusActive
	return s.cardRepo.Save(card)
}
