package travel

import (
	"errors"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/user"
)

func TestUseCase_UpdateTravelPlan(t *testing.T) {
	userID := uint(1)
	travelID := uint(1)
	travelPlanID := uint(1)

	date := time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)
	newDate := time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC)

	place := "バンコク"
	newPlace := "チェンマイ"

	content := "ワット・ポー"
	newContent := "ワット・プラ・シン"

	sortOrder := 1
	newSortOrder := 2

	t.Run("正常系", func(t *testing.T) {
		tests := []struct {
			name        string
			updateInput UpdateTravelPlanInput
			want        *TravelPlan
		}{
			{
				name: "日付を更新する",
				updateInput: UpdateTravelPlanInput{
					StartDate: &newDate,
				},
				want: &TravelPlan{
					ID:        travelPlanID,
					TravelID:  travelID,
					StartDate: newDate,
					Place:     place,
					Content:   content,
					SortOrder: sortOrder,
				},
			},
			{
				name: "場所を更新する",
				updateInput: UpdateTravelPlanInput{
					Place: &newPlace,
				},
				want: &TravelPlan{
					ID:        travelPlanID,
					TravelID:  travelID,
					StartDate: date,
					Place:     newPlace,
					Content:   content,
					SortOrder: sortOrder,
				},
			},
			{
				name: "内容を更新する",
				updateInput: UpdateTravelPlanInput{
					Content: &newContent,
				},
				want: &TravelPlan{
					ID:        travelPlanID,
					TravelID:  travelID,
					StartDate: date,
					Place:     place,
					Content:   newContent,
					SortOrder: sortOrder,
				},
			},
			{
				name: "順番を更新する",
				updateInput: UpdateTravelPlanInput{
					SortOrder: &newSortOrder,
				},
				want: &TravelPlan{
					ID:        travelPlanID,
					TravelID:  travelID,
					StartDate: date,
					Place:     place,
					Content:   content,
					SortOrder: newSortOrder,
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repository := &MockPlanRepository{
					FindByIDFunc: func(id uint) (*TravelPlan, error) {
						return &TravelPlan{
							ID:        travelPlanID,
							TravelID:  travelID,
							StartDate: date,
							Place:     place,
							Content:   content,
							SortOrder: sortOrder,
						}, nil
					},
					UpdateFunc: func(travelPlan *TravelPlan) error {
						return nil
					},
				}

				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*Travel, error) {
						return &Travel{
							ID:     travelID,
							UserID: userID,
						}, nil
					},
				}

				useCase := NewPlanUseCase(travelRepository, repository)

				got, err := useCase.UpdateTravelPlan(
					userID,
					travelID,
					travelPlanID,
					tt.updateInput,
				)

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("expected %+v, got %+v", tt.want, got)
				}
			})
		}
	})

	t.Run("異常系", func(t *testing.T) {
		databaseError := errors.New("database error")
		updateError := errors.New("update error")

		tests := []struct {
			name                    string
			travelRepository        *Travel
			travelRepositoryErr     error
			travelPlanRepository    *TravelPlan
			travelPlanRepositoryErr error
			updateErr               error
			wantErr                 error
		}{
			{
				name:                "Travelが存在しない",
				travelRepositoryErr: gorm.ErrRecordNotFound,
				wantErr:             ErrTravelNotFound,
			},
			{
				name: "他人のTravelを更新しようとする",
				travelRepository: &Travel{
					ID:     travelID,
					UserID: 2,
				},
				wantErr: user.ErrUserNotAuthorized,
			},
			{
				name:                "TravelのRepositoryでエラーが発生する",
				travelRepositoryErr: databaseError,
				wantErr:             databaseError,
			},
			{
				name: "TravelPlanが存在しない",
				travelRepository: &Travel{
					ID:     travelID,
					UserID: userID,
				},
				travelPlanRepositoryErr: gorm.ErrRecordNotFound,
				wantErr:                 ErrTravelPlanNotFound,
			},
			{
				name: "TravelPlanのRepositoryでエラーが発生する",
				travelRepository: &Travel{
					ID:     travelID,
					UserID: userID,
				},
				travelPlanRepositoryErr: databaseError,
				wantErr:                 databaseError,
			},
			{
				name: "TravelPlanが指定されたTravelに属していない",
				travelRepository: &Travel{
					ID:     travelID,
					UserID: userID,
				},
				travelPlanRepository: &TravelPlan{
					ID:       travelPlanID,
					TravelID: 2,
				},
				wantErr: ErrTravelPlanNotBelongToTravel,
			},
			{
				name: "TravelPlanのUpdateでエラーが発生する",
				travelRepository: &Travel{
					ID:     travelID,
					UserID: userID,
				},
				travelPlanRepository: &TravelPlan{
					ID:       travelPlanID,
					TravelID: travelID,
				},
				updateErr: updateError,
				wantErr:   updateError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repository := &MockPlanRepository{
					FindByIDFunc: func(id uint) (*TravelPlan, error) {
						return tt.travelPlanRepository, tt.travelPlanRepositoryErr
					},
					UpdateFunc: func(travelPlan *TravelPlan) error {
						return tt.updateErr
					},
				}

				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*Travel, error) {
						return tt.travelRepository, tt.travelRepositoryErr
					},
				}

				useCase := NewPlanUseCase(travelRepository, repository)

				_, err := useCase.UpdateTravelPlan(
					userID,
					travelID,
					travelPlanID,
					UpdateTravelPlanInput{
						Place: &newPlace,
					},
				)

				if !errors.Is(err, tt.wantErr) {
					t.Fatalf(
						"expected error %v, got %v",
						tt.wantErr,
						err,
					)
				}
			})
		}
	})
}
