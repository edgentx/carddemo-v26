package dto

// OpenAccountRequest defines the JSON payload for creating an account.
type OpenAccountRequest struct {
	UserProfileID string `json:"user_profile_id" validate:"required"`
	Status        string `json:"status" validate:"required"`
	AccountType   string `json:"account_type" validate:"required"`
}

// UpdateAccountStatusRequest defines the JSON payload for updating status.
type UpdateAccountStatusRequest struct {
	NewStatus string `json:"new_status" validate:"required"`
	Reason    string `json:"reason"`
}

// AccountResponse defines the JSON response for an account.
type AccountResponse struct {
	ID        string `json:"id"`
	ProfileID string `json:"profile_id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
