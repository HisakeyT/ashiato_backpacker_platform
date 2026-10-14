package user

import (
	"errors"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

var _ Repository = (*GormRepository)(nil)

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *GormRepository) FindByID(id uint) (*model.User, error) {
	var user model.User

	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *GormRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *GormRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}
