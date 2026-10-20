package usecase

import (
	"net/url"
	"strings"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
)

type PlanRecordUsecase struct {
	travels TravelRepository
	plans   repository.TravelPlanRepository
	records repository.PlanRecordRepository
}

func NewPlanRecordUsecase(travels TravelRepository, plans repository.TravelPlanRepository, records repository.PlanRecordRepository) *PlanRecordUsecase {
	return &PlanRecordUsecase{travels: travels, plans: plans, records: records}
}

type RecordInput struct {
	Name         string
	MealType     *string
	LocationText *string
	MapURL       *string
}

type RecordPatch struct {
	Name         *string
	MealType     **string
	LocationText **string
	MapURL       **string
}

func validMapURL(value *string) bool {
	if value == nil {
		return true
	}
	u, err := url.Parse(*value)
	return err == nil && (u.Scheme == "https" || u.Scheme == "http") && u.Host != "" && u.User == nil
}

func validMealType(value *string) bool {
	if value == nil {
		return true
	}
	switch *value {
	case "breakfast", "lunch", "dinner", "snack", "other":
		return true
	default:
		return false
	}
}

func validateRecord(kind model.RecordKind, record *model.PlanRecord) error {
	if !kind.Valid() || strings.TrimSpace(record.Name) == "" || !validMapURL(record.MapURL) {
		return model.ErrInvalidPlanRecord
	}
	if kind == model.RecordMeal {
		if !validMealType(record.MealType) {
			return model.ErrInvalidPlanRecord
		}
	} else if record.MealType != nil {
		return model.ErrInvalidPlanRecord
	}
	return nil
}

func (uc *PlanRecordUsecase) plan(userID, travelID, planID uint, public bool) error {
	travel, err := uc.travels.FindByID(travelID)
	if err != nil {
		return err
	}
	if public {
		if !travel.IsPublic {
			return model.ErrTravelNotFound
		}
	} else if travel.UserID != userID {
		return model.ErrUserNotAuthorized
	}
	plan, err := uc.plans.FindByID(planID)
	if err != nil {
		return err
	}
	if plan.TravelID != travelID {
		return model.ErrTravelPlanNotFound
	}
	return nil
}

func (uc *PlanRecordUsecase) List(userID, travelID, planID uint, kind model.RecordKind, public bool) ([]model.PlanRecord, error) {
	if !kind.Valid() {
		return nil, model.ErrInvalidPlanRecord
	}
	if err := uc.plan(userID, travelID, planID, public); err != nil {
		return nil, err
	}
	return uc.records.List(kind, planID)
}

func (uc *PlanRecordUsecase) Create(userID, travelID, planID uint, kind model.RecordKind, input RecordInput) (*model.PlanRecord, error) {
	if err := uc.plan(userID, travelID, planID, false); err != nil {
		return nil, err
	}
	record := &model.PlanRecord{TravelPlanID: planID, Name: strings.TrimSpace(input.Name), MealType: input.MealType, LocationText: input.LocationText, MapURL: input.MapURL}
	if err := validateRecord(kind, record); err != nil {
		return nil, err
	}
	if err := uc.records.Create(kind, record); err != nil {
		return nil, err
	}
	return record, nil
}

func (uc *PlanRecordUsecase) Update(userID, travelID, planID, id uint, kind model.RecordKind, patch RecordPatch) (*model.PlanRecord, error) {
	if err := uc.plan(userID, travelID, planID, false); err != nil {
		return nil, err
	}
	record, err := uc.records.Find(kind, id)
	if err != nil {
		return nil, err
	}
	if record.TravelPlanID != planID {
		return nil, model.ErrPlanRecordNotFound
	}
	if patch.Name != nil {
		record.Name = strings.TrimSpace(*patch.Name)
	}
	if patch.MealType != nil {
		record.MealType = *patch.MealType
	}
	if patch.LocationText != nil {
		record.LocationText = *patch.LocationText
	}
	if patch.MapURL != nil {
		record.MapURL = *patch.MapURL
	}
	if err := validateRecord(kind, record); err != nil {
		return nil, err
	}
	if err := uc.records.Update(kind, record); err != nil {
		return nil, err
	}
	return record, nil
}

func (uc *PlanRecordUsecase) Delete(userID, travelID, planID, id uint, kind model.RecordKind) error {
	if err := uc.plan(userID, travelID, planID, false); err != nil {
		return err
	}
	record, err := uc.records.Find(kind, id)
	if err != nil {
		return err
	}
	if record.TravelPlanID != planID {
		return model.ErrPlanRecordNotFound
	}
	return uc.records.Delete(kind, id)
}
