package command

// ReportCardLostCmd is a command to flag the card as lost.
// This disables the card immediately and triggers a replacement process.
type ReportCardLostCmd struct {
	CardID      string
	LossReason  string
	ReportedBy  string // ID of the user or system reporting the loss
	// ForceUsage is a testing hook to force a specific daily usage value
	// without needing to construct a complex transaction history.
	ForceUsage *int
	// ForceStatus is a testing hook to set a specific initial status
	// for the purpose of testing state transition invariants.
	ForceStatus string
}
