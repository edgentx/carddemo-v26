package event

// TransactionSubmitted is emitted when a transaction is successfully validated and created.
type TransactionSubmitted struct {
	TransactionID   string
	AccountID       string
	CardID          string
	Amount          float64
	TransactionType string
}

// TransactionReversed is emitted when a transaction is successfully reversed.
type TransactionReversed struct {
	TransactionID string
	Reason        string
}
