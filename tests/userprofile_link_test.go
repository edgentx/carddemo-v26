package tests

import (
	"testing"

	"github.com/carddemo/project/src/domain/userprofile/command"
	"github.com/carddemo/project/src/domain/userprofile/model"
	"github.com/carddemo/project/src/domain/userprofile/repository"
	"github.com/carddemo/project/mocks"
)

func TestLinkUserToAccountCmd_Success(t *testing.T) {
	// Setup mock repo
	repo := mocks.NewMockUserProfileRepository()

	// Given a valid UserProfile aggregate
	user := model.NewUserProfile("user-123", repo)

	// And a valid user_id is provided (already set in constructor)
	// And a valid account_id is provided
	cmd := command.LinkUserToAccountCmd{
		AggregateID: "user-123",
		AccountID:   "acc-456",
		IsVerified:  true,
	}

	// When the LinkUserToAccountCmd command is executed
	evts, err := user.Execute(cmd)

	// Then no error is returned
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// And a user.linked.to.account event is emitted
	if len(evts) != 1 {
		t.Fatalf("expected 1 event, got %d", len(evts))
	}

	evt := evts[0]
	if evt.Type() != "com.carddemo.user.linked.to.account" {
		t.Errorf("expected event type 'com.carddemo.user.linked.to.account', got '%s'", evt.Type())
	}

	// And state is updated
	if user.AccountID != "acc-456" {
		t.Errorf("expected AccountID 'acc-456', got '%s'", user.AccountID)
	}
	if !user.IsPrimary {
		t.Errorf("expected IsPrimary to be true")
	}
}

func TestLinkUserToAccountCmd_Rejected_NotVerified(t *testing.T) {
	// Setup mock repo
	repo := mocks.NewMockUserProfileRepository()

	// Given a UserProfile aggregate
	user := model.NewUserProfile("user-123", repo)

	// And a command that violates: A user profile must contain verified identity information
	cmd := command.LinkUserToAccountCmd{
		AggregateID: "user-123",
		AccountID:   "acc-456",
		IsVerified:  false,
	}

	// When the LinkUserToAccountCmd command is executed
	evts, err := user.Execute(cmd)

	// Then the command is rejected with a domain error
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err != model.ErrIdentityNotVerified {
		t.Errorf("expected ErrIdentityNotVerified, got %v", err)
	}

	if len(evts) != 0 {
		t.Errorf("expected no events, got %d", len(evts))
	}
}

func TestLinkUserToAccountCmd_Rejected_AlreadyPrimary(t *testing.T) {
	// Setup mock repo
	repo := mocks.NewMockUserProfileRepository()

	// Given a UserProfile aggregate that violates: An active account must have exactly one primary user profile
	// We simulate this by saving another user profile as primary for the same account to the mock repo.
	existingUser := model.NewUserProfile("user-999", nil) // nil repo for this one to avoid recursion
	existingUser.AccountID = "acc-456"
	existingUser.IsPrimary = true
	repo.Save(existingUser)

	// Now we try to link a new user to the same account
	user := model.NewUserProfile("user-123", repo)

	cmd := command.LinkUserToAccountCmd{
		AggregateID: "user-123",
		AccountID:   "acc-456",
		IsVerified:  true,
	}

	// When the LinkUserToAccountCmd command is executed
	evts, err := user.Execute(cmd)

	// Then the command is rejected with a domain error
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err != model.ErrDuplicatePrimaryProfile {
		t.Errorf("expected ErrDuplicatePrimaryProfile, got %v", err)
	}

	if len(evts) != 0 {
		t.Errorf("expected no events, got %d", len(evts))
	}
}
