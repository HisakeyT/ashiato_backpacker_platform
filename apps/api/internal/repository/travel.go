package repository

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

type TravelRepository interface {
	Create(travel *model.Travel) error
	FindByID(id uint) (*model.Travel, error)
	FindByUserID(userID uint) ([]*model.Travel, error)
	Update(travel *model.Travel) error
	FindPublicTravels() ([]*model.PublicTravel, error)
	FindPublicByID(id uint) (*model.PublicTravel, error)
	FindCopySource(sourceTravelID uint) (*model.CopySource, error)
	Delete(id uint) error
	CopyTravel(originalTravel *model.Travel, newUserID uint) (*model.Travel, error)
}
