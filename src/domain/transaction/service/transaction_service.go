package service

import (
	"github.com/carddemo/project/src/domain/transaction/command"
	"github.com/carddemo/project/src/domain/transaction/model"
)

// TransactionService defines domain services that act on transactions.
// This is used if logic needs to be pulled out of the aggregate, but for this story,
// validation can live inside the aggregate or use a simple interface.
// We define this to adhere to the src/domain/{module}/service/ requirement if needed.

// ValidateTransactionForSubmission simulates complex domain logic if necessary.
// For TDD, we invoke the aggregate directly.
func ValidateTransactionForSubmission(t *model.Transaction, cmd command.SubmitTransactionCmd) error {
	// Placeholder for domain service logic if needed.
	// In this implementation, the aggregate handles the logic.
	return nil
}
