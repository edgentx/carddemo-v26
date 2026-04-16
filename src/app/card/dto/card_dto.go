package dto

// IssueCardRequest represents the request payload for issuing a card
type IssueCardRequest struct {
	AccountID      string         `json:"account_id"`
	CardType       string         `json:"card_type"`
	SpendingLimits map[string]int `json:"spending_limits"`
}

// UpdateCardStatusRequest represents the request payload for updating status
type UpdateCardStatusRequest struct {
	Status string `json:"status"`
}

// UpdatePolicyRequest represents the request payload for updating policy
type UpdatePolicyRequest struct {
	DailyLimit int `json:"daily_limit"`
}
