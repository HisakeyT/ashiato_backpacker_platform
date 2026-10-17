package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
)

func TestTravelPlanTextValidation(t *testing.T) {
	startDate := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	for _, operation := range []string{"create", "update"} {
		for _, field := range []string{"place", "content"} {
			t.Run(operation+"/"+field, func(t *testing.T) {
				written := false
				travelRepository := &mock.MockTravelRepository{
					FindByIDFunc: func(id uint) (*model.Travel, error) { return &model.Travel{ID: 1, UserID: 1}, nil },
				}
				planRepository := &mock.MockTravelPlanRepository{
					FindByIDFunc: func(id uint) (*model.TravelPlan, error) {
						return &model.TravelPlan{ID: 2, TravelID: 1, StartDate: startDate, Place: "place", Content: "content"}, nil
					},
					CreateFunc: func(plan *model.TravelPlan) error { written = true; return nil },
					UpdateFunc: func(plan *model.TravelPlan) error { written = true; return nil },
				}
				useCase := NewTravelPlanUsecase(travelRepository, planRepository)
				place, content := "place", "content"
				if field == "place" {
					place = "  "
				} else {
					content = "  "
				}
				var err error
				if operation == "create" {
					_, err = useCase.CreateTravelPlan(1, 1, CreateTravelPlanInput{StartDate: startDate, Place: place, Content: content})
				} else {
					_, err = useCase.UpdateTravelPlan(1, 1, 2, UpdateTravelPlanInput{Place: &place, Content: &content})
				}
				if !errors.Is(err, model.ErrInvalidTravelPlanText) || written {
					t.Fatalf("error = %v, written = %v", err, written)
				}
			})
		}
	}
}
