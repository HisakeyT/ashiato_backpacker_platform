package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
)

var errRepository = errors.New("repository error")

func TestUsecase_CreateTravelPlan(t *testing.T) {
	t.Run("正常に作成できる", func(t *testing.T) {
		// Arrange
		userID := uint(1)
		travelID := uint(10)

		input := CreateTravelPlanInput{
			StartDate: planDate(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)),
			Place:     "バンコク",
			Content:   "ワット・ポーを観光",
			SortOrder: 1,
		}

		travelData := &model.Travel{
			ID:        travelID,
			UserID:    userID,
			StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC),
		}

		travelRepository := &mock.MockTravelRepository{
			FindByIDFunc: func(id uint) (*model.Travel, error) {
				return travelData, nil
			},
		}

		travelPlanRepository := &mock.MockTravelPlanRepository{
			CreateFunc: func(travelPlan *model.TravelPlan) error {
				return nil
			},
		}

		useCase := NewTravelPlanUsecase(
			travelRepository,
			travelPlanRepository,
		)

		// Act
		travelPlan, err := useCase.CreateTravelPlan(userID, travelID, input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if travelPlan.TravelID != travelID {
			t.Errorf("TravelID = %d, want %d", travelPlan.TravelID, travelID)
		}

		if travelPlan.StartDate == nil || !travelPlan.StartDate.Equal(*input.StartDate) {
			t.Errorf("StartDate = %v, want %v", travelPlan.StartDate, input.StartDate)
		}

		if travelPlan.Place != input.Place {
			t.Errorf("Place = %s, want %s", travelPlan.Place, input.Place)
		}

		if travelPlan.Content != input.Content {
			t.Errorf("Content = %s, want %s", travelPlan.Content, input.Content)
		}

		if travelPlan.SortOrder != input.SortOrder {
			t.Errorf("SortOrder = %d, want %d", travelPlan.SortOrder, input.SortOrder)
		}
	})
}

func TestUsecase_CreateTravelPlan_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		travelID    uint
		travelData  *model.Travel
		findByIDErr error
		wantErr     error
	}{
		{
			name:     "他人の旅行には旅程を追加できない",
			userID:   1,
			travelID: 10,
			travelData: &model.Travel{
				ID:     10,
				UserID: 2,
			},
			wantErr: model.ErrUserNotAuthorized,
		},
		{
			name:        "旅行の取得に失敗した場合",
			userID:      1,
			travelID:    10,
			findByIDErr: errRepository,
			wantErr:     errRepository,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			travelRepository := &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return tt.travelData, tt.findByIDErr
				},
			}

			travelPlanRepository := &mock.MockTravelPlanRepository{
				CreateFunc: func(travelPlan *model.TravelPlan) error {
					t.Fatal("Create should not be called")
					return nil
				},
			}

			useCase := NewTravelPlanUsecase(
				travelRepository,
				travelPlanRepository,
			)

			input := CreateTravelPlanInput{
				StartDate: planDate(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)),
				Place:     "バンコク",
				Content:   "観光",
				SortOrder: 1,
			}

			_, err := useCase.CreateTravelPlan(
				tt.userID,
				tt.travelID,
				input,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsecase_CreateTravelPlan_Location(t *testing.T) {
	country := "  Japan  "
	city := "   "

	tests := []struct {
		name     string
		location model.Location
		wantErr  error
	}{
		{
			name: "所在地なしで作成できる",
		},
		{
			name: "所在地を正規化して保存する",
			location: model.Location{
				Country:   &country,
				City:      &city,
				Latitude:  new(35.0),
				Longitude: new(135.0),
			},
		},
		{
			name: "緯度だけの指定はエラー",
			location: model.Location{
				Latitude: new(35.0),
			},
			wantErr: model.ErrInvalidCoordinates,
		},
		{
			name: "緯度が範囲外",
			location: model.Location{
				Latitude:  new(91.0),
				Longitude: new(135.0),
			},
			wantErr: model.ErrInvalidCoordinates,
		},
		{
			name: "経度が範囲外",
			location: model.Location{
				Latitude:  new(35.0),
				Longitude: new(181.0),
			},
			wantErr: model.ErrInvalidCoordinates,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			travelRepository := &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:        10,
						UserID:    1,
						StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC),
					}, nil
				},
			}

			var saved *model.TravelPlan
			travelPlanRepository := &mock.MockTravelPlanRepository{
				CreateFunc: func(plan *model.TravelPlan) error {
					saved = plan
					return nil
				},
			}

			uc := NewTravelPlanUsecase(
				travelRepository,
				travelPlanRepository,
			)

			got, err := uc.CreateTravelPlan(
				1,
				10,
				CreateTravelPlanInput{
					StartDate: planDate(time.Date(
						2026, 10, 1, 0, 0, 0, 0, time.UTC,
					)),
					Place:    "京都",
					Content:  "観光",
					Location: tt.location,
				},
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr != nil {
				if got != nil {
					t.Fatalf("expected nil result, got %+v", got)
				}
				if saved != nil {
					t.Fatal("Create must not be called for invalid location")
				}
				return
			}

			if got == nil || saved == nil {
				t.Fatal("expected a created travel plan")
			}

			if tt.location.Country != nil {
				if saved.Location.Country == nil ||
					*saved.Location.Country != "Japan" {
					t.Fatalf(
						"expected normalized country Japan, got %v",
						saved.Location.Country,
					)
				}
			}

			if saved.Location.City != nil {
				t.Fatalf(
					"expected nil city, got %q",
					*saved.Location.City,
				)
			}

			if tt.location.Latitude != nil {
				if saved.Location.Latitude == nil ||
					*saved.Location.Latitude != *tt.location.Latitude {
					t.Fatal("latitude must be preserved")
				}
			}

			if tt.location.Longitude != nil {
				if saved.Location.Longitude == nil ||
					*saved.Location.Longitude != *tt.location.Longitude {
					t.Fatal("longitude must be preserved")
				}
			}
		})
	}
}
