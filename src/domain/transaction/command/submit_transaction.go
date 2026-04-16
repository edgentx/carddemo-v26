package command

// SubmitTransactionCmd is the command to submit a new financial movement.
type SubmitTransactionCmd struct {
	TransactionID   string
	AccountID       string
	CardID          string
	Amount          float64
	TransactionType string // e.g., "debit", "credit"
	AccountStatus   string // "Active", "Inactive", etc. (Simulated fetch)
}
