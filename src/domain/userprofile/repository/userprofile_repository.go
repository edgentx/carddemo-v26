package repository

import (
	"github.com/carddemo/project/src/domain/userprofile/model"
)

// UserProfileRepository defines the persistence interface for UserProfiles.
type UserProfileRepository interface {
	Get(id string) (*model.UserProfile, error)
	Save(aggregate *model.UserProfile) error
	Delete(id string) error
}
