package userprofile

import "errors"

var (
	// ErrIdentityNotVerified is returned when identification details are not verified.
	ErrIdentityNotVerified = errors.New("user profile must contain verified identity information to be linked to an account")

	// ErrDuplicatePrimaryProfile is returned when an active account already has a primary profile.
	ErrDuplicatePrimaryProfile = errors.New("an active account must have exactly one primary user profile")
)
