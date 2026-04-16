package model

import (
	"errors"

	"github.com/carddemo/project/src/domain/shared"
	"github.com/carddemo/project/src/domain/userprofile/command"
	"github.com/carddemo/project/src/domain/userprofile/event"
	"github.com/carddemo/project/src/domain/userprofile/repository"
)

var (
	// ErrIdentityNotVerified is returned when identification details are not verified.
	ErrIdentityNotVerified = errors.New("user profile must contain verified identity information to be linked to an account")

	// ErrDuplicatePrimaryProfile is returned when an active account already has a primary profile.
	ErrDuplicatePrimaryProfile = errors.New("an active account must have exactly one primary user profile")
)

// UserProfile represents the UserProfile aggregate.
// It maintains its own state and enforces invariants upon command execution.
type UserProfile struct {
	shared.AggregateRoot
	ID           string
	Email        string
	CreditScore  int
	IsPrimary    bool
	IsRegistered bool
	AccountID    string

	// repo is injected to allow the Aggregate to verify cross-aggregate invariants
	// such as uniqueness of a primary profile within an account context.
	repo repository.UserProfileRepository
}

// NewUserProfile creates a new UserProfile instance.
// It accepts an optional repository to enable invariant checking.
func NewUserProfile(id string, repo ...repository.UserProfileRepository) *UserProfile {
	var r repository.UserProfileRepository
	if len(repo) > 0 {
		r = repo[0]
	}

	return &UserProfile{
		ID:   id,
		repo: r,
	}
}

// InjectRepository allows setting the repository after creation if needed.
func (u *UserProfile) InjectRepository(repo repository.UserProfileRepository) {
	u.repo = repo
}

// Execute handles commands for the UserProfile aggregate.
func (u *UserProfile) Execute(cmd interface{}) ([]shared.DomainEvent, error) {
	switch c := cmd.(type) {
	case command.RegisterUserCmd:
		return u.handleRegisterUser(c)
	case command.LinkUserToAccountCmd:
		return u.handleLinkUserToAccount(c)
	default:
		return nil, shared.ErrUnknownCommand
	}
}

// handleRegisterUser processes the RegisterUserCmd command.
func (u *UserProfile) handleRegisterUser(cmd command.RegisterUserCmd) ([]shared.DomainEvent, error) {
	// Invariant 1: Identity must be verified.
	if !cmd.Identification.IsVerified {
		return nil, ErrIdentityNotVerified
	}

	// Logic to check for existing primary profiles using the repo if available.
	if u.repo != nil {
		existing, _ := u.repo.List()
		for _, ex := range existing {
			if ex.ID != u.ID && ex.IsPrimary {
				return nil, ErrDuplicatePrimaryProfile
			}
		}
	}

	u.Email = cmd.ContactInfo.Email
	u.CreditScore = cmd.CreditProfile.Score
	u.IsRegistered = true

	evt := event.NewUserRegistered(u.ID, u.Email, u.CreditScore)
	return []shared.DomainEvent{evt}, nil
}

// handleLinkUserToAccount processes the LinkUserToAccountCmd command.
func (u *UserProfile) handleLinkUserToAccount(cmd command.LinkUserToAccountCmd) ([]shared.DomainEvent, error) {
	// Invariant: Identity must be verified to link to an account.
	if !cmd.IsVerified {
		return nil, ErrIdentityNotVerified
	}

	// Invariant: Check if another user is already primary for this account.
	if u.repo != nil {
		existing, _ := u.repo.List()
		for _, ex := range existing {
			if ex.ID != u.ID && ex.AccountID == cmd.AccountID && ex.IsPrimary {
				return nil, ErrDuplicatePrimaryProfile
			}
		}
	}

	// Apply state changes
	u.AccountID = cmd.AccountID
	u.IsPrimary = true

	evt := event.NewUserLinkedToAccount(u.ID, u.AccountID)
	return []shared.DomainEvent{evt}, nil
}
