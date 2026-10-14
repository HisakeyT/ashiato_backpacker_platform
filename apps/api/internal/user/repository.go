package user

import (
	"errors"
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type Repository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
}
