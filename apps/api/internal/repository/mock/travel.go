package mock

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

type MockTravelRepository struct {
	CreateFunc            func(travel *model.Travel) error
	FindByIDFunc          func(id uint) (*model.Travel, error)
	FindByUserIDFunc      func(userID uint) ([]*model.Travel, error)
	UpdateFunc            func(travel *model.Travel) error
	FindPublicTravelsFunc func() ([]*model.PublicTravel, error)
	FindPublicByIDFunc    func(id uint) (*model.PublicTravel, error)
	FindCopySourceFunc    func(sourceTravelID uint) (*model.CopySource, error)
	DeleteFunc            func(id uint) error
	FindByIDWithPlansFunc func(id uint) (*model.Travel, error)
	CreateWithPlansFunc   func(src *model.Travel) error
}

func (m *MockTravelRepository) Create(travel *model.Travel) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(travel)
	}
	return nil
}

func (m *MockTravelRepository) FindByID(id uint) (*model.Travel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}

	return nil, nil
}

func (m *MockTravelRepository) FindByUserID(userID uint) ([]*model.Travel, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(userID)
	}

	return nil, nil
}

func (m *MockTravelRepository) Update(travel *model.Travel) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(travel)
	}

	return nil
}

func (m *MockTravelRepository) FindPublicTravels() ([]*model.PublicTravel, error) {
	if m.FindPublicTravelsFunc != nil {
		return m.FindPublicTravelsFunc()
	}

	return nil, nil
}

func (m *MockTravelRepository) FindPublicByID(id uint) (*model.PublicTravel, error) {
	if m.FindPublicByIDFunc != nil {
		return m.FindPublicByIDFunc(id)
	}

	return nil, nil
}

func (m *MockTravelRepository) FindCopySource(sourceTravelID uint) (*model.CopySource, error) {
	if m.FindCopySourceFunc != nil {
		return m.FindCopySourceFunc(sourceTravelID)
	}

	return nil, nil
}

func (m *MockTravelRepository) Delete(id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockTravelRepository) FindByIDWithPlans(id uint) (*model.Travel, error) {
	if m.FindByIDWithPlansFunc != nil {
		return m.FindByIDWithPlansFunc(id)
	}

	return nil, nil
}

func (m *MockTravelRepository) CreateWithPlans(src *model.Travel) error {
	if m.CreateWithPlansFunc != nil {
		return m.CreateWithPlansFunc(src)
	}

	return nil
}
