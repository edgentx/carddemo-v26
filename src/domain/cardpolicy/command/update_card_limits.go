package command

// UpdateCardLimitsCmd is a command to adjust the daily or monthly spending limits
// allocated to a specific card.
type UpdateCardLimitsCmd struct {
	CardID             string
	DailyLimit         int64
	MonthlyLimit       int64
	AuthorizationToken string
	ProfileID          string // Used to trigger invariants based on risk profile
}
