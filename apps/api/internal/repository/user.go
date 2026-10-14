package repository

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

type Repository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
}
