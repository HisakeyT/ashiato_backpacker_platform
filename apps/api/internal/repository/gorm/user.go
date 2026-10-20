package gormrepo

import (
	"errors"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User

	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Model(&model.User{}).Where("id = ?", user.ID).Update("nickname", user.Nickname).Error
}

// BumpAuthVersion invalidates all existing tokens for this user. Credential updates
// can call this on a repository bound to the same transaction before commit.
func (r *UserRepository) BumpAuthVersion(id uint) error {
	result := r.db.Model(&model.User{}).Where("id = ?", id).
		UpdateColumn("auth_version", gorm.Expr("auth_version + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

var _ repository.UserRepository = (*UserRepository)(nil)
var _ repository.AuthVersionRepository = (*UserRepository)(nil)
