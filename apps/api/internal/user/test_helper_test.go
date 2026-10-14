package user

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

type MockRepository struct {
	CreateFunc      func(user *model.User) error
	FindByIDFunc    func(id uint) (*model.User, error)
	FindByEmailFunc func(email string) (*model.User, error)
	UpdateFunc      func(user *model.User) error
}

func (m *MockRepository) Create(user *model.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

func (m *MockRepository) FindByID(id uint) (*model.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, ErrUserNotFound
}

func (m *MockRepository) FindByEmail(email string) (*model.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(email)
	}
	return nil, ErrUserNotFound
}

func (m *MockRepository) Update(user *model.User) error {
	return m.UpdateFunc(user)
}
