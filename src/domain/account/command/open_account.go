package command

// OpenAccountCmd is a command to create a new account.
type OpenAccountCmd struct {
	UserProfileID  string
	InitialStatus  string
	AccountType    string
}
