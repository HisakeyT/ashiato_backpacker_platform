package travel

import (
	"errors"
	"github.com/HisakeyT/backpacker-platform/internal/user"
	"time"

	"gorm.io/gorm"
)

type TravelRepository interface {
	FindByID(id uint) (*Travel, error)
}

type PlanUseCase struct {
	travelRepository     TravelRepository
	travelPlanRepository PlanRepository
}

func NewPlanUseCase(
	travelRepository TravelRepository,
	travelPlanRepository PlanRepository,
) *PlanUseCase {
	return &PlanUseCase{
		travelRepository:     travelRepository,
		travelPlanRepository: travelPlanRepository,
	}
}

type CreateTravelPlanInput struct {
	StartDate time.Time
	Place     string
	Content   string
	SortOrder int
}

func (uc *PlanUseCase) CreateTravelPlan(userID, travelID uint, input CreateTravelPlanInput) (*TravelPlan, error) {
	travelData, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTravelNotFound
		}
		return nil, err
	}

	if travelData.UserID != userID {
		return nil, user.ErrUserNotAuthorized
	}
	travelPlan := &TravelPlan{
		TravelID:  travelID,
		StartDate: input.StartDate,
		Place:     input.Place,
		Content:   input.Content,
		SortOrder: input.SortOrder,
	}

	if err := uc.travelPlanRepository.Create(travelPlan); err != nil {
		return nil, err
	}

	return travelPlan, nil
}

func (uc *PlanUseCase) GetTravelPlans(userID, travelID uint) ([]TravelPlan, error) {
	travelData, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTravelNotFound
		}
		return nil, err
	}

	if travelData.UserID != userID {
		return nil, user.ErrUserNotAuthorized
	}

	travelPlans, err := uc.travelPlanRepository.FindByTravelID(travelID)
	if err != nil {
		return nil, err
	}

	return travelPlans, nil
}

func (uc *PlanUseCase) DeleteTravelPlan(userID, travelID, travelPlanID uint) error {
	travelData, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTravelNotFound
		}
		return err
	}

	if travelData.UserID != userID {
		return user.ErrUserNotAuthorized
	}

	travelPlan, err := uc.travelPlanRepository.FindByID(travelPlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTravelPlanNotFound
		}
		return err
	}

	if travelPlan.TravelID != travelID {
		return ErrTravelPlanNotBelongToTravel
	}

	return uc.travelPlanRepository.Delete(travelPlanID)
}

type UpdateTravelPlanInput struct {
	StartDate *time.Time
	Place     *string
	Content   *string
	SortOrder *int
}

func (uc *PlanUseCase) UpdateTravelPlan(userID, travelID, travelPlanID uint, input UpdateTravelPlanInput) (*TravelPlan, error) {
	travelData, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTravelNotFound
		}
		return nil, err
	}

	if travelData.UserID != userID {
		return nil, user.ErrUserNotAuthorized
	}

	travelPlan, err := uc.travelPlanRepository.FindByID(travelPlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTravelPlanNotFound
		}

		return nil, err
	}

	if travelPlan.TravelID != travelID {
		return nil, ErrTravelPlanNotBelongToTravel
	}

	applyTravelPlanUpdates(travelPlan, input)

	if err := uc.travelPlanRepository.Update(travelPlan); err != nil {
		return nil, err
	}

	return travelPlan, nil
}

func applyTravelPlanUpdates(travelPlan *TravelPlan, input UpdateTravelPlanInput) {
	if input.StartDate != nil {
		travelPlan.StartDate = *input.StartDate
	}

	if input.Place != nil {
		travelPlan.Place = *input.Place
	}

	if input.Content != nil {
		travelPlan.Content = *input.Content
	}

	if input.SortOrder != nil {
		travelPlan.SortOrder = *input.SortOrder
	}
}

func (uc *PlanUseCase) GetPulicTravelPlan(travelID uint) ([]TravelPlan, error) {
	travelData, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTravelNotFound
		}
		return nil, err
	}

	if !travelData.IsPublic {
		return nil, ErrTravelNotFound
	}

	travelPlans, err := uc.travelPlanRepository.FindByTravelID(travelID)
	if err != nil {
		return nil, err
	}

	return travelPlans, nil
}
