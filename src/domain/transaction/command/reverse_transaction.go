package command

// ReverseTransactionCmd is the command to reverse a previously submitted transaction.
type ReverseTransactionCmd struct {
	TransactionID string
	Amount        float64
	AccountStatus string
	Reason        string
}
