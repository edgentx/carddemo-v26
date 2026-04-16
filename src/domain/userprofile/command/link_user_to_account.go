package command

// LinkUserToAccountCmd is the command to link a user profile to a financial account.
type LinkUserToAccountCmd struct {
	AggregateID string
	AccountID   string
	IsVerified  bool // Flag indicating if identity verification is complete
}
