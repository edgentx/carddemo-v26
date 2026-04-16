package event

// TransactionSubmitted is emitted when a transaction is successfully validated and created.
// This acts as a fact in the system.
type TransactionSubmitted struct {
	TransactionID   string
	AccountID       string
	CardID          string
	Amount          float64
	TransactionType string
}
