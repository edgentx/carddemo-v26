package repository

import "github.com/carddemo/project/src/domain/userprofile/model"

// UserProfileRepository defines the storage interface for UserProfile aggregates.
type UserProfileRepository interface {
	Get(id string) (*model.UserProfile, error)
	Save(aggregate *model.UserProfile) error
	Delete(id string) error
	List() ([]*model.UserProfile, error)
}
