package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Example of wrapping tests in a suite if preferred for BDD style setup/teardown

type DomainTestSuite struct {
	suite.Suite
	Repo    *MockAccountRepository
	Service *AccountDomainService
}

func (suite *DomainTestSuite) SetupTest() {
	suite.Repo = NewMockAccountRepository()
	suite.Service = NewAccountDomainService(suite.Repo)
}

func (suite *DomainTestSuite) TearDownTest() {
	// Cleanup
}

func TestDomainTestSuite(t *testing.T) {
	suite.Run(t, new(DomainTestSuite))
}
