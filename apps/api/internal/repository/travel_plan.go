package repository

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

type TravelPlanRepository interface {
	Create(travelPlan *model.TravelPlan) error
	Update(travelPlan *model.TravelPlan) error
	Delete(id uint) error
	FindByID(id uint) (*model.TravelPlan, error)
	FindByTravelID(travelID uint) ([]model.TravelPlan, error)
	Reorder(travelID uint, ids []uint) error
}
