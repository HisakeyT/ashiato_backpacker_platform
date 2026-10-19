package usecase

import (
	"strings"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
)

type TravelRepository interface {
	FindByID(id uint) (*model.Travel, error)
}

type TravelPlanUsecase struct {
	travelRepository     TravelRepository
	travelPlanRepository repository.TravelPlanRepository
}

func NewTravelPlanUsecase(
	travelRepository TravelRepository,
	travelPlanRepository repository.TravelPlanRepository,
) *TravelPlanUsecase {
	return &TravelPlanUsecase{
		travelRepository:     travelRepository,
		travelPlanRepository: travelPlanRepository,
	}
}

type CreateTravelPlanInput struct {
	StartDate *time.Time
	EndDate   *time.Time
	Place     string
	Location  model.Location
	Content   string
	SortOrder int
}

func (uc *TravelPlanUsecase) CreateTravelPlan(userID, travelID uint, input CreateTravelPlanInput) (*model.TravelPlan, error) {
	travelData, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		return nil, err
	}

	if travelData.UserID != userID {
		return nil, model.ErrUserNotAuthorized
	}

	if err := validatePlanDates(input.StartDate, input.EndDate, travelData); err != nil {
		return nil, err
	}

	input.Location.Normalize()

	if err := input.Location.Validate(); err != nil {
		return nil, err
	}
	if err := validatePlanText(input.Place, input.Content); err != nil {
		return nil, err
	}

	travelPlan := &model.TravelPlan{
		TravelID:  travelID,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		Location:  input.Location,
		Place:     input.Place,
		Content:   input.Content,
		SortOrder: input.SortOrder,
	}

	if err := uc.travelPlanRepository.Create(travelPlan); err != nil {
		return nil, err
	}

	return travelPlan, nil
}

type UpdateTravelPlanInput struct {
	StartDate    *time.Time
	StartDateSet bool
	EndDate      *time.Time
	EndDateSet   bool
	Place        *string
	Location     *model.Location
	Content      *string
	SortOrder    *int
}

func (uc *TravelPlanUsecase) UpdateTravelPlan(userID, travelID, travelPlanID uint, input UpdateTravelPlanInput) (*model.TravelPlan, error) {
	travelData, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		return nil, err
	}

	if travelData.UserID != userID {
		return nil, model.ErrUserNotAuthorized
	}

	travelPlan, err := uc.travelPlanRepository.FindByID(travelPlanID)
	if err != nil {
		return nil, err
	}

	if travelPlan.TravelID != travelID {
		return nil, model.ErrTravelPlanNotBelongToTravel
	}

	applyTravelPlanUpdates(travelPlan, input)

	if err := validatePlanDates(travelPlan.StartDate, travelPlan.EndDate, travelData); err != nil {
		return nil, err
	}

	if err := travelPlan.Location.Validate(); err != nil {
		return nil, err
	}
	if err := validatePlanText(travelPlan.Place, travelPlan.Content); err != nil {
		return nil, err
	}

	if err := uc.travelPlanRepository.Update(travelPlan); err != nil {
		return nil, err
	}

	return travelPlan, nil
}

func validatePlanDates(start *time.Time, end *time.Time, travel *model.Travel) error {
	if start == nil {
		if end != nil {
			return model.ErrTravelPlanEndWithoutStart
		}
		if travel.ScheduleMode == model.ScheduleModeRecorded {
			return model.ErrTravelPlanStartDateRequired
		}
		return nil
	}
	lastDay := *start
	if end != nil {
		lastDay = *end
	}
	if err := validateDateRange(*start, lastDay); err != nil {
		return err
	}
	if start.Before(travel.StartDate) || lastDay.After(travel.EndDate) {
		return model.ErrTravelPlanOutsideTravelDates
	}
	return nil
}

func validatePlanText(place, content string) error {
	if strings.TrimSpace(place) == "" || strings.TrimSpace(content) == "" {
		return model.ErrInvalidTravelPlanText
	}
	return nil
}

func (uc *TravelPlanUsecase) GetTravelPlans(userID, travelID uint) ([]model.TravelPlan, error) {
	travelData, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		return nil, err
	}

	if travelData.UserID != userID {
		return nil, model.ErrUserNotAuthorized
	}

	travelPlans, err := uc.travelPlanRepository.FindByTravelID(travelID)
	if err != nil {
		return nil, err
	}

	return travelPlans, nil
}

func (uc *TravelPlanUsecase) DeleteTravelPlan(userID, travelID, travelPlanID uint) error {
	travelData, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		return err
	}

	if travelData.UserID != userID {
		return model.ErrUserNotAuthorized
	}

	travelPlan, err := uc.travelPlanRepository.FindByID(travelPlanID)
	if err != nil {
		return err
	}

	if travelPlan.TravelID != travelID {
		return model.ErrTravelPlanNotBelongToTravel
	}

	return uc.travelPlanRepository.Delete(travelPlanID)
}

func applyTravelPlanUpdates(travelPlan *model.TravelPlan, input UpdateTravelPlanInput) {
	if input.StartDateSet || input.StartDate != nil {
		travelPlan.StartDate = input.StartDate
	}

	if input.EndDateSet || input.EndDate != nil {
		travelPlan.EndDate = input.EndDate
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

	if input.Location != nil {
		travelPlan.Location = *input.Location
		travelPlan.Location.Normalize()
	}
}

func (uc *TravelPlanUsecase) GetPublicTravelPlan(travelID uint) ([]model.TravelPlan, error) {
	travelData, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		return nil, err
	}

	if !travelData.IsPublic {
		return nil, model.ErrTravelNotFound
	}

	travelPlans, err := uc.travelPlanRepository.FindByTravelID(travelID)
	if err != nil {
		return nil, err
	}

	return travelPlans, nil
}

func (uc *TravelPlanUsecase) ReorderTravelPlans(userID, travelID uint, ids []uint) error {
	travel, err := uc.travelRepository.FindByID(travelID)
	if err != nil {
		return err
	}
	if travel.UserID != userID {
		return model.ErrUserNotAuthorized
	}
	if travel.ScheduleMode != model.ScheduleModePlanning {
		return model.ErrTravelPlanOrderOnlyPlanning
	}
	return uc.travelPlanRepository.Reorder(travelID, ids)
}
