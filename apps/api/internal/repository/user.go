package repository

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
}

// AuthVersionRepository is used by credential updates inside their transaction.
type AuthVersionRepository interface {
	BumpAuthVersion(id uint) error
}
