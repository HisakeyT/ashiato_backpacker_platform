package gormrepo

import (
	"errors"

	"gorm.io/gorm"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
)

type TravelRepository struct {
	db *gorm.DB
}

func NewTravelRepository(db *gorm.DB) *TravelRepository {
	return &TravelRepository{db: db}
}

func (r *TravelRepository) Create(travel *model.Travel) error {
	return r.db.Create(travel).Error
}

func (r *TravelRepository) FindByID(id uint) (*model.Travel, error) {
	var travel model.Travel

	if err := r.db.First(&travel, id).Error; err != nil {
		return nil, err
	}

	return &travel, nil
}

func (r *TravelRepository) FindByUserID(userID uint) ([]*model.Travel, error) {
	var travels []*model.Travel

	if err := r.db.Where("user_id = ?", userID).Find(&travels).Error; err != nil {
		return nil, err
	}

	return travels, nil
}

func (r *TravelRepository) Update(travel *model.Travel) error {
	return r.db.Save(travel).Error
}

func (r *TravelRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Travel{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *TravelRepository) FindPublicTravels() ([]*model.PublicTravel, error) {
	var PublicTravels []*model.PublicTravel

	if err := r.db.Model(&model.Travel{}).
		Select("travels.*, users.nickname AS author_nickname").
		Joins("JOIN users ON users.id = travels.user_id").
		Where("travels.is_public = ?", true).
		Order("travels.created_at DESC").
		Scan(&PublicTravels).Error; err != nil {
		return nil, err
	}

	return PublicTravels, nil
}

func (r *TravelRepository) FindPublicByID(id uint) (*model.PublicTravel, error) {
	var publicTravel model.PublicTravel

	result := r.db.Model(&model.Travel{}).
		Select("travels.*, users.nickname AS author_nickname").
		Joins("JOIN users ON users.id = travels.user_id").
		Where("travels.is_public = ? AND travels.id = ?", true, id).
		Limit(1).
		Scan(&publicTravel)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &publicTravel, nil
}

func (r *TravelRepository) FindCopySource(sourceTravelID uint) (*model.CopySource, error) {
	var copySource model.CopySource

	result := r.db.Model(&model.Travel{}).
		Select("travels.id AS travel_id, users.nickname AS author_nickname").
		Joins("JOIN users ON users.id = travels.user_id").
		Where("travels.is_public = ? AND travels.id = ?", true, sourceTravelID).
		Limit(1).
		Scan(&copySource)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &copySource, nil
}

func (r *TravelRepository) FindByIDWithPlans(id uint) (*model.Travel, error) {
	var travel model.Travel
	err := r.db.
		Preload("Plans", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, id ASC")
		}).
		First(&travel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrTravelNotFound
		}
		return nil, err
	}
	return &travel, nil
}

func (r *TravelRepository) CreateWithPlans(travel *model.Travel) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Plans").Create(travel).Error; err != nil {
			return err
		}
		if len(travel.Plans) == 0 {
			return nil
		}
		for i := range travel.Plans {
			travel.Plans[i].TravelID = travel.ID
		}
		return tx.Create(&travel.Plans).Error
	})
}

var _ repository.TravelRepository = (*TravelRepository)(nil)
