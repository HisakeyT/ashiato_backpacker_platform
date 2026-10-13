package travel

import (
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
)

func TestUseCase_CreateTravel(t *testing.T) {
	userID := uint(1)

	startDate := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		input CreateTravelInput
	}{
		{
			name: "正常に作成できる",
			input: CreateTravelInput{
				Title:     "東南アジア3週間",
				StartDate: startDate,
				EndDate:   endDate,
				IsPublic:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &MockRepository{
				CreateFunc: func(travel *model.Travel) error {
					travel.ID = 1
					return nil
				},
			}

			useCase := NewUseCase(repository)

			travel, err := useCase.CreateTravel(userID, tt.input)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if travel.ID != 1 {
				t.Fatalf("expected ID 1, got %d", travel.ID)
			}

			if travel.UserID != userID {
				t.Fatalf("expected user ID 1, got %d", travel.UserID)
			}

			if travel.Title != tt.input.Title {
				t.Fatalf("expected title %s, got %s", tt.input.Title, travel.Title)
			}

			if !travel.StartDate.Equal(tt.input.StartDate) {
				t.Fatalf("unexpected start date")
			}

			if !travel.EndDate.Equal(tt.input.EndDate) {
				t.Fatalf("unexpected end date")
			}

			if travel.IsPublic != tt.input.IsPublic {
				t.Fatalf("expected IsPublic %v, got %v", tt.input.IsPublic, travel.IsPublic)
			}
		})
	}
}

func TestUseCase_CreateTravel_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input CreateTravelInput
	}{
		{
			name: "開始日が終了日より後",
			input: CreateTravelInput{
				Title:     "東南アジア3週間",
				StartDate: time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
				IsPublic:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &MockRepository{}
			useCase := NewUseCase(repository)

			_, err := useCase.CreateTravel(1, tt.input)

			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
