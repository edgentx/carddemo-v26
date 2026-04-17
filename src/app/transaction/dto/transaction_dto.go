package dto

// CreateTransactionRequest defines the JSON payload for creating a transaction.
type CreateTransactionRequest struct {
	AccountID       string  `json:"account_id" validate:"required"`
	CardID          string  `json:"card_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=debit credit"`
}

// VoidTransactionRequest defines the JSON payload for voiding a transaction.
type VoidTransactionRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// TransactionResponse defines the JSON response for a transaction.
type TransactionResponse struct {
	ID              string  `json:"id"`
	AccountID       string  `json:"account_id"`
	CardID          string  `json:"card_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	Status          string  `json:"status"`
	CreatedAt       string  `json:"created_at"`
}

// BatchSettlementRequest defines the JSON payload for creating a batch settlement.
type BatchSettlementRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// BatchSettlementResponse defines the JSON response for a batch settlement.
type BatchSettlementResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}
