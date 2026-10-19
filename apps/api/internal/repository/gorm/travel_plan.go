package gormrepo

import (
	"errors"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormPlanRepository struct {
	db *gorm.DB
}

func NewGormPlanRepository(db *gorm.DB) *GormPlanRepository {
	return &GormPlanRepository{db: db}
}

func (r *GormPlanRepository) Create(travelPlan *model.TravelPlan) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var travel model.Travel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&travel, travelPlan.TravelID).Error; err != nil {
			return err
		}
		var maxOrder int
		if err := tx.Model(&model.TravelPlan{}).Where("travel_id = ?", travelPlan.TravelID).
			Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder).Error; err != nil {
			return err
		}
		travelPlan.SortOrder = maxOrder + 1
		return tx.Create(travelPlan).Error
	})
}

func (r *GormPlanRepository) Update(travelPlan *model.TravelPlan) error {
	return r.db.Save(travelPlan).Error
}

func (r *GormPlanRepository) Delete(id uint) error {
	return r.db.Delete(&model.TravelPlan{}, id).Error
}

func (r *GormPlanRepository) FindByID(id uint) (*model.TravelPlan, error) {
	var travelPlan model.TravelPlan

	result := r.db.First(&travelPlan, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, model.ErrTravelPlanNotFound
		}
		return nil, result.Error
	}

	return &travelPlan, nil
}

func (r *GormPlanRepository) FindByTravelID(travelID uint) ([]model.TravelPlan, error) {
	var travelPlans []model.TravelPlan

	result := r.db.Table("travel_plans").
		Joins("JOIN travels ON travels.id = travel_plans.travel_id").
		Where("travel_plans.travel_id = ?", travelID).
		Order("CASE WHEN travels.schedule_mode = 'planning' THEN travel_plans.sort_order END ASC").
		Order("CASE WHEN travels.schedule_mode = 'recorded' THEN travel_plans.start_date END ASC").
		Order("CASE WHEN travels.schedule_mode = 'recorded' THEN travel_plans.end_date END ASC").
		Order("CASE WHEN travels.schedule_mode = 'recorded' THEN travel_plans.sort_order END ASC").
		Order("travel_plans.id ASC").Select("travel_plans.*").Find(&travelPlans)
	if result.Error != nil {
		return nil, result.Error
	}

	return travelPlans, nil
}

func (r *GormPlanRepository) Reorder(travelID uint, ids []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var travel model.Travel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&travel, travelID).Error; err != nil {
			return err
		}
		if travel.ScheduleMode != model.ScheduleModePlanning {
			return model.ErrTravelPlanOrderOnlyPlanning
		}
		var plans []model.TravelPlan
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("travel_id = ?", travelID).Find(&plans).Error; err != nil {
			return err
		}
		currentIDs := make([]uint, len(plans))
		for i, plan := range plans {
			currentIDs[i] = plan.ID
		}
		if err := validatePlanOrder(currentIDs, ids); err != nil {
			return err
		}
		for i, id := range ids {
			if err := tx.Model(&model.TravelPlan{}).Where("id = ? AND travel_id = ?", id, travelID).Update("sort_order", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func validatePlanOrder(currentIDs, requestedIDs []uint) error {
	if len(currentIDs) != len(requestedIDs) {
		return model.ErrInvalidTravelPlanOrder
	}
	known := make(map[uint]bool, len(currentIDs))
	for _, id := range currentIDs {
		known[id] = true
	}
	for _, id := range requestedIDs {
		if !known[id] {
			return model.ErrInvalidTravelPlanOrder
		}
		delete(known, id)
	}
	return nil
}
