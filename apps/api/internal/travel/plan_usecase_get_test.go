package travel

import (
	"errors"
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
)

func TestUseCase_GetTravelPlans(t *testing.T) {
	userID := uint(1)
	travelID := uint(10)

	t.Run("正常系", func(t *testing.T) {
		tests := []struct {
			name          string
			travelPlans   []model.TravelPlan
			expectedCount int
		}{
			{
				name: "TravelPlanを複数件取得できる",
				travelPlans: []model.TravelPlan{
					{
						ID:        1,
						TravelID:  travelID,
						StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
						Place:     "バンコク",
						Content:   "ワット・ポーを観光",
						SortOrder: 1,
					},
					{
						ID:        2,
						TravelID:  travelID,
						StartDate: time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC),
						Place:     "バンコク",
						Content:   "王宮を観光",
						SortOrder: 2,
					},
				},
				expectedCount: 2,
			},
			{
				name:          "TravelPlanが0件でも取得できる",
				travelPlans:   []model.TravelPlan{},
				expectedCount: 0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				travelData := &model.Travel{
					ID:     travelID,
					UserID: userID,
				}

				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*model.Travel, error) {
						return travelData, nil
					},
				}

				travelPlanRepository := &MockPlanRepository{
					FindByTravelIDFunc: func(id uint) ([]model.TravelPlan, error) {
						return tt.travelPlans, nil
					},
				}

				useCase := NewPlanUseCase(
					travelRepository,
					travelPlanRepository,
				)

				// Act
				travelPlans, err := useCase.GetTravelPlans(userID, travelID)

				// Assert
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if len(travelPlans) != tt.expectedCount {
					t.Errorf(
						"len(travelPlans) = %d, want %d",
						len(travelPlans),
						tt.expectedCount,
					)
				}
			})
		}
	})

	t.Run("異常系", func(t *testing.T) {
		tests := []struct {
			name          string
			travelData    *model.Travel
			findByIDError error
			expectedError error
		}{
			{
				name:          "Travelが存在しない",
				findByIDError: gorm.ErrRecordNotFound,
				expectedError: model.ErrTravelNotFound,
			},
			{
				name: "自分のTravelではない",
				travelData: &model.Travel{
					ID:     travelID,
					UserID: 999,
				},
				expectedError: model.ErrUserNotAuthorized,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*model.Travel, error) {
						return tt.travelData, tt.findByIDError
					},
				}

				travelPlanRepository := &MockPlanRepository{}

				useCase := NewPlanUseCase(
					travelRepository,
					travelPlanRepository,
				)

				// Act
				_, err := useCase.GetTravelPlans(userID, travelID)

				// Assert
				if !errors.Is(err, tt.expectedError) {
					t.Errorf(
						"error = %v, want %v",
						err,
						tt.expectedError,
					)
				}
			})
		}
	})
}
