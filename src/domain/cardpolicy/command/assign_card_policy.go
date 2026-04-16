package command

// AssignCardPolicyCmd is a command to assign a specific set of transaction rules
// and restrictions to a card.
type AssignCardPolicyCmd struct {
	CardID               string
	PolicyType           string
	MerchantRestrictions []string
}
