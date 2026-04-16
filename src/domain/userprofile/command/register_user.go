package command

// RegisterUserCmd is the command to register a new user.
type RegisterUserCmd struct {
	AggregateID string
	ContactInfo    ContactInfoVO
	CreditProfile  CreditProfileVO
	Identification IdentificationVO
}

type ContactInfoVO struct {
	Email string
	Phone string
}

type CreditProfileVO struct {
	Score int
}

type IdentificationVO struct {
	IDNumber string
	IsVerified bool
}
