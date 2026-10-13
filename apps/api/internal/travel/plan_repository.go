package travel

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

type PlanRepository interface {
	Create(travelPlan *model.TravelPlan) error
	Update(travelPlan *model.TravelPlan) error
	Delete(id uint) error
	FindByID(id uint) (*model.TravelPlan, error)
	FindByTravelID(travelID uint) ([]model.TravelPlan, error)
}
