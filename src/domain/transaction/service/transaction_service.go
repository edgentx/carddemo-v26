package service

// Commands for Transaction Aggregate

type SubmitTransactionCmd struct {
	TransactionID   string
	AccountID       string
	CardID          string
	Amount          float64
	TransactionType string
	AccountStatus   string
}

type ReverseTransactionCmd struct {
	TransactionID string
	Amount        float64
	AccountStatus string
	Reason        string
}
