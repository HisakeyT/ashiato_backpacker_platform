package travel

import (
	"gorm.io/gorm"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(travel *model.Travel) error {
	return r.db.Create(travel).Error
}

func (r *GormRepository) FindByID(id uint) (*model.Travel, error) {
	var travel model.Travel

	if err := r.db.First(&travel, id).Error; err != nil {
		return nil, err
	}

	return &travel, nil
}

func (r *GormRepository) FindByUserID(userID uint) ([]*model.Travel, error) {
	var travels []*model.Travel

	if err := r.db.Where("user_id = ?", userID).Find(&travels).Error; err != nil {
		return nil, err
	}

	return travels, nil
}

func (r *GormRepository) Update(travel *model.Travel) error {
	return r.db.Save(travel).Error
}

func (r *GormRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Travel{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormRepository) FindPublicTravels() ([]*model.PublicTravel, error) {
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

func (r *GormRepository) FindPublicByID(id uint) (*model.PublicTravel, error) {
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

func (r *GormRepository) FindCopySource(sourceTravelID uint) (*model.CopySource, error) {
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

func (r *GormRepository) CopyTravel(originalTravel *model.Travel, newUserID uint) (*model.Travel, error) {
	var createdTravel model.Travel

	err := r.db.Transaction(func(tx *gorm.DB) error {
		createdTravel = model.Travel{
			UserID:             newUserID,
			Title:              originalTravel.Title,
			StartDate:          originalTravel.StartDate,
			EndDate:            originalTravel.EndDate,
			CopiedFromTravelID: &originalTravel.ID,
		}
		if err := tx.Create(&createdTravel).Error; err != nil {
			return err
		}

		var travelPlans []model.TravelPlan
		if err := tx.Where("travel_id = ?", originalTravel.ID).Find(&travelPlans).Error; err != nil {
			return err
		}
		if len(travelPlans) == 0 {
			return nil
		}

		copiedPlans := make([]model.TravelPlan, len(travelPlans))
		for i, plan := range travelPlans {
			var endDate *time.Time
			if plan.EndDate != nil {
				e := *plan.EndDate
				endDate = &e
			}
			copiedPlans[i] = model.TravelPlan{
				TravelID:  createdTravel.ID,
				StartDate: plan.StartDate,
				EndDate:   endDate,
				Place:     plan.Place,
				Location:  plan.Location,
				Content:   plan.Content,
				SortOrder: plan.SortOrder,
			}
		}
		return tx.Create(&copiedPlans).Error
	})
	if err != nil {
		return nil, err
	}

	return &createdTravel, nil
}
