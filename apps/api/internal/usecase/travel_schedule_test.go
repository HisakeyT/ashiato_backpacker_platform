package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
)

func TestPlanDatesByScheduleMode(t *testing.T) {
	start := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	travel := &model.Travel{StartDate: start, EndDate: start.AddDate(0, 0, 3), ScheduleMode: model.ScheduleModePlanning}
	if err := validatePlanDates(nil, nil, travel); err != nil {
		t.Fatalf("undated planning plan: %v", err)
	}
	end := start.AddDate(0, 0, 1)
	if err := validatePlanDates(nil, &end, travel); !errors.Is(err, model.ErrTravelPlanEndWithoutStart) {
		t.Fatalf("end only: %v", err)
	}
	travel.ScheduleMode = model.ScheduleModeRecorded
	if err := validatePlanDates(nil, nil, travel); !errors.Is(err, model.ErrTravelPlanStartDateRequired) {
		t.Fatalf("recorded plan without start: %v", err)
	}
}

func TestUpdateTravelRecordingNeedsAllPlanStartDates(t *testing.T) {
	start := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	mode := model.ScheduleModeRecorded
	travel := &model.Travel{ID: 1, UserID: 2, StartDate: start, EndDate: start.AddDate(0, 0, 3), ScheduleMode: model.ScheduleModePlanning}
	updated := false
	repo := &mock.MockTravelRepository{
		FindByIDFunc: func(uint) (*model.Travel, error) { return travel, nil },
		FindByIDWithPlansFunc: func(uint) (*model.Travel, error) {
			return &model.Travel{Plans: []model.TravelPlan{{ID: 5, StartDate: nil}}}, nil
		},
		UpdateFunc: func(*model.Travel) error { updated = true; return nil },
	}
	_, err := NewTravelUsecase(repo).UpdateTravel(2, 1, UpdateTravelInput{ScheduleMode: &mode})
	if !errors.Is(err, model.ErrPlanStartDateRequired) || updated {
		t.Fatalf("error = %v, updated = %v", err, updated)
	}
}

func TestReorderPlansRequiresOwnerAndPlanningMode(t *testing.T) {
	travel := &model.Travel{ID: 1, UserID: 2, ScheduleMode: model.ScheduleModePlanning}
	called := false
	uc := NewTravelPlanUsecase(
		&mock.MockTravelRepository{FindByIDFunc: func(uint) (*model.Travel, error) { return travel, nil }},
		&mock.MockTravelPlanRepository{ReorderFunc: func(uint, []uint) error { called = true; return nil }},
	)
	if err := uc.ReorderTravelPlans(3, 1, []uint{1}); !errors.Is(err, model.ErrUserNotAuthorized) {
		t.Fatalf("other user: %v", err)
	}
	travel.ScheduleMode = model.ScheduleModeRecorded
	if err := uc.ReorderTravelPlans(2, 1, []uint{1}); !errors.Is(err, model.ErrTravelPlanOrderOnlyPlanning) {
		t.Fatalf("recorded: %v", err)
	}
	if called {
		t.Fatal("repository called for forbidden reorder")
	}
	travel.ScheduleMode = model.ScheduleModePlanning
	if err := uc.ReorderTravelPlans(2, 1, []uint{1}); err != nil || !called {
		t.Fatalf("planning reorder: %v, called = %v", err, called)
	}
}
