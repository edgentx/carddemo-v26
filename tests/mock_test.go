package tests

import (
	"testing"

	"github.com/carddemo/project/mocks"
	"github.com/carddemo/project/src/domain/account/repository"
	"github.com/carddemo/project/src/domain/account/model"
)

// TestMockRepositoriesImplementInterfaces ensures mocks adhere to domain contracts.
func TestMockRepositoriesImplementInterfaces(t *testing.T) {
	// Compile-time check for Account Repository
	var _ repository.AccountRepository = (*mocks.MockAccountRepository)(nil)

	// Runtime check
	repo := mocks.NewMockAccountRepository()
	if repo == nil {
		t.Fatal("Failed to instantiate MockAccountRepository")
	}

	// Functional check
	acc := model.NewAccount("999")
	err := repo.Save(acc)
	if err != nil {
		t.Errorf("Mock Save failed: %v", err)
	}

	retrieved, err := repo.Get("999")
	if err != nil {
		t.Errorf("Mock Get failed: %v", err)
	}
	if retrieved.ID != "999" {
		t.Errorf("Mock retrieved wrong ID: %s", retrieved.ID)
	}
}
