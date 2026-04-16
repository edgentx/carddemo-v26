package command

// IssueCardCmd is a command to request the creation of a new card.
type IssueCardCmd struct {
	AccountID       string
	CardType        string
	SpendingLimits  map[string]int
	IsLostOrStolen  bool // Used to force state for testing invariants
	DailyTxnLimit   int  // Used to force state for testing invariants
	CurrentUsage    int  // Used to force state for testing invariants
}
