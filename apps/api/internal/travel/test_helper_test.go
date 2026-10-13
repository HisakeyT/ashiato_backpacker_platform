package travel

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

type MockRepository struct {
	CreateFunc            func(travel *model.Travel) error
	FindByIDFunc          func(id uint) (*model.Travel, error)
	FindByUserIDFunc      func(userID uint) ([]*model.Travel, error)
	UpdateFunc            func(travel *model.Travel) error
	FindPublicTravelsFunc func() ([]*model.PublicTravel, error)
	FindPublicByIDFunc    func(id uint) (*model.PublicTravel, error)
	FindCopySourceFunc    func(sourceTravelID uint) (*model.CopySource, error)
	DeleteFunc            func(id uint) error
	CopyTravelFunc        func(originalTravel *model.Travel, newUserID uint) (*model.Travel, error)
}

func (m *MockRepository) Create(travel *model.Travel) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(travel)
	}
	return nil
}

func (m *MockRepository) FindByID(id uint) (*model.Travel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}

	return nil, nil
}

func (m *MockRepository) FindByUserID(userID uint) ([]*model.Travel, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(userID)
	}

	return nil, nil
}

func (m *MockRepository) Update(travel *model.Travel) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(travel)
	}

	return nil
}

func (m *MockRepository) FindPublicTravels() ([]*model.PublicTravel, error) {
	if m.FindPublicTravelsFunc != nil {
		return m.FindPublicTravelsFunc()
	}

	return nil, nil
}

func (m *MockRepository) FindPublicByID(id uint) (*model.PublicTravel, error) {
	if m.FindPublicByIDFunc != nil {
		return m.FindPublicByIDFunc(id)
	}

	return nil, nil
}

func (m *MockRepository) FindCopySource(sourceTravelID uint) (*model.CopySource, error) {
	if m.FindCopySourceFunc != nil {
		return m.FindCopySourceFunc(sourceTravelID)
	}

	return nil, nil
}

func (m *MockRepository) Delete(id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockRepository) CopyTravel(originalTravel *model.Travel, newUserID uint) (*model.Travel, error) {
	if m.CopyTravelFunc != nil {
		return m.CopyTravelFunc(originalTravel, newUserID)
	}

	return nil, nil
}
