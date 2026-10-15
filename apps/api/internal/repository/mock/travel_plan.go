package mock

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

type MockTravelPlanRepository struct {
	CreateFunc         func(travelPlan *model.TravelPlan) error
	UpdateFunc         func(travelPlan *model.TravelPlan) error
	DeleteFunc         func(id uint) error
	FindByIDFunc       func(id uint) (*model.TravelPlan, error)
	FindByTravelIDFunc func(travelID uint) ([]model.TravelPlan, error)
}

func (m *MockTravelPlanRepository) Create(travelPlan *model.TravelPlan) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(travelPlan)
	}
	return nil
}

func (m *MockTravelPlanRepository) Update(travelPlan *model.TravelPlan) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(travelPlan)
	}

	return nil
}

func (m *MockTravelPlanRepository) Delete(id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}

	return nil
}

func (m *MockTravelPlanRepository) FindByID(id uint) (*model.TravelPlan, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}

	return nil, nil
}

func (m *MockTravelPlanRepository) FindByTravelID(travelID uint) ([]model.TravelPlan, error) {
	if m.FindByTravelIDFunc != nil {
		return m.FindByTravelIDFunc(travelID)
	}

	return nil, nil
}
