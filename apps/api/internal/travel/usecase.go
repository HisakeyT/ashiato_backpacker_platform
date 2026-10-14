package travel

import (
	"errors"
	"gorm.io/gorm"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
)

var ErrInvalidDateRange = errors.New("start date must not be after end date")

type UseCase struct {
	repository Repository
}

func NewUseCase(repository Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

type CreateTravelInput struct {
	Title     string
	StartDate time.Time
	EndDate   time.Time
	IsPublic  bool
}

func (u *UseCase) CreateTravel(userID uint, input CreateTravelInput) (*model.Travel, error) {
	if err := validateDateRange(input.StartDate, input.EndDate); err != nil {
		return nil, err
	}

	travel := &model.Travel{
		UserID:    userID,
		Title:     input.Title,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		IsPublic:  input.IsPublic,
	}

	if err := u.repository.Create(travel); err != nil {
		return nil, err
	}

	return travel, nil

}

func validateDateRange(startDate, endDate time.Time) error {
	if startDate.After(endDate) {
		return ErrInvalidDateRange
	}

	return nil
}

func (u *UseCase) GetPublicTravels() ([]*model.PublicTravel, error) {
	return u.repository.FindPublicTravels()
}

// ロジックがないので、テストは書かない
func (u *UseCase) GetPublicTravel(travelID uint) (*model.PublicTravel, error) {
	return u.repository.FindPublicByID(travelID)
}

func (u *UseCase) GetTravels(userID uint) ([]*model.Travel, error) {
	travels, err := u.repository.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	return travels, nil
}

func (u *UseCase) GetTravel(userID uint, travelID uint) (*model.TravelDetail, error) {
	travel, err := u.repository.FindByID(travelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrTravelNotFound
		}
		return nil, err
	}

	if travel.UserID != userID {
		return nil, model.ErrUserNotAuthorized
	}

	var copiedFrom *model.CopySource
	if travel.CopiedFromTravelID != nil {
		copiedFrom, err = u.repository.FindCopySource(*travel.CopiedFromTravelID)
		if err != nil {
			return nil, err
		}
	}

	travelDetail := &model.TravelDetail{Travel: travel, CopiedFrom: copiedFrom}
	return travelDetail, nil
}

type UpdateTravelInput struct {
	Title     *string    `json:"title"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsPublic  *bool      `json:"is_public"`
}

func (u *UseCase) UpdateTravel(userID uint, travelID uint, input UpdateTravelInput) (*model.Travel, error) {
	travel, err := u.repository.FindByID(travelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrTravelNotFound
		}
		return nil, err
	}

	if travel.UserID != userID {
		return nil, model.ErrUserNotAuthorized
	}

	applyUpdate(travel, input)

	if err := validateDateRange(travel.StartDate, travel.EndDate); err != nil {
		return nil, err
	}

	if err := u.repository.Update(travel); err != nil {
		return nil, err
	}

	return travel, nil
}

func applyUpdate(travel *model.Travel, input UpdateTravelInput) {
	if input.Title != nil {
		travel.Title = *input.Title
	}

	if input.StartDate != nil {
		travel.StartDate = *input.StartDate
	}

	if input.EndDate != nil {
		travel.EndDate = *input.EndDate
	}

	if input.IsPublic != nil {
		travel.IsPublic = *input.IsPublic
	}
}

func (u *UseCase) DeleteTravel(userID uint, travelID uint) error {
	travel, err := u.repository.FindByID(travelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.ErrTravelNotFound
		}

		return err
	}

	if travel.UserID != userID {
		return model.ErrUserNotAuthorized
	}

	if err := u.repository.Delete(travelID); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) CopyTravel(userID, travelID uint) (*model.Travel, error) {
	originalTravel, err := uc.repository.FindByID(travelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrTravelNotFound
		}
		return nil, err
	}

	// 非公開かつ他人の旅行は存在しない扱い
	if !originalTravel.IsPublic && originalTravel.UserID != userID {
		return nil, model.ErrTravelNotFound
	}

	return uc.repository.CopyTravel(originalTravel, userID)
}
