package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
)

func TestTravelPlanDatesWithinTravel(t *testing.T) {
	date := func(day int) time.Time { return time.Date(2026, 10, day, 0, 0, 0, 0, time.UTC) }
	datePtr := func(day int) *time.Time { value := date(day); return &value }
	travel := &model.Travel{ID: 1, UserID: 1, StartDate: date(1), EndDate: date(3)}
	tests := []struct {
		name  string
		start int
		end   *int
		want  error
	}{
		{"first day without end", 1, nil, nil},
		{"last day without end", 3, nil, nil},
		{"full travel range", 1, new(3), nil},
		{"same day end", 2, new(2), nil},
		{"before travel", 0, nil, model.ErrTravelPlanOutsideTravelDates},
		{"after travel without end", 4, nil, model.ErrTravelPlanOutsideTravelDates},
		{"end after travel", 2, new(4), model.ErrTravelPlanOutsideTravelDates},
		{"end before start", 2, new(1), model.ErrInvalidDateRange},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var end *time.Time
			if tt.end != nil {
				end = datePtr(*tt.end)
			}
			for _, operation := range []string{"create", "update"} {
				t.Run(operation, func(t *testing.T) {
					written := false
					travelRepository := &mock.MockTravelRepository{FindByIDFunc: func(uint) (*model.Travel, error) { return travel, nil }}
					planRepository := &mock.MockTravelPlanRepository{
						FindByIDFunc: func(uint) (*model.TravelPlan, error) {
							return &model.TravelPlan{ID: 2, TravelID: 1, StartDate: date(1), Place: "place", Content: "content"}, nil
						},
						CreateFunc: func(*model.TravelPlan) error { written = true; return nil },
						UpdateFunc: func(*model.TravelPlan) error { written = true; return nil },
					}
					uc := NewTravelPlanUsecase(travelRepository, planRepository)
					var err error
					if operation == "create" {
						_, err = uc.CreateTravelPlan(1, 1, CreateTravelPlanInput{StartDate: date(tt.start), EndDate: end, Place: "place", Content: "content"})
					} else {
						_, err = uc.UpdateTravelPlan(1, 1, 2, UpdateTravelPlanInput{StartDate: datePtr(tt.start), EndDate: end, EndDateSet: true})
					}
					if !errors.Is(err, tt.want) {
						t.Fatalf("error = %v, want %v", err, tt.want)
					}
					if written != (tt.want == nil) {
						t.Fatalf("written = %v, want %v", written, tt.want == nil)
					}
				})
			}
		})
	}
}
