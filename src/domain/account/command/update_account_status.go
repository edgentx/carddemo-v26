package command

// UpdateAccountStatusCmd is a command to change the lifecycle state of an account.
type UpdateAccountStatusCmd struct {
	NewStatus string
	Reason    string
}
