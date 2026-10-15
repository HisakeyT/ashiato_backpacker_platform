package usecase

import (
	"errors"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
)

func TestUsecase_UpdateTravelPlan(t *testing.T) {
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
			want        *model.TravelPlan
		}{
			{
				name: "日付を更新する",
				updateInput: UpdateTravelPlanInput{
					StartDate: &newDate,
				},
				want: &model.TravelPlan{
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
				want: &model.TravelPlan{
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
				want: &model.TravelPlan{
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
				want: &model.TravelPlan{
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
				repository := &mock.MockTravelPlanRepository{
					FindByIDFunc: func(id uint) (*model.TravelPlan, error) {
						return &model.TravelPlan{
							ID:        travelPlanID,
							TravelID:  travelID,
							StartDate: date,
							Place:     place,
							Content:   content,
							SortOrder: sortOrder,
						}, nil
					},
					UpdateFunc: func(travelPlan *model.TravelPlan) error {
						return nil
					},
				}

				travelRepository := &mock.MockTravelRepository{
					FindByIDFunc: func(id uint) (*model.Travel, error) {
						return &model.Travel{
							ID:     travelID,
							UserID: userID,
						}, nil
					},
				}

				useCase := NewTravelPlanUsecase(travelRepository, repository)

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
			travelRepository        *model.Travel
			travelRepositoryErr     error
			travelPlanRepository    *model.TravelPlan
			travelPlanRepositoryErr error
			updateErr               error
			wantErr                 error
		}{
			{
				name:                "Travelが存在しない",
				travelRepositoryErr: gorm.ErrRecordNotFound,
				wantErr:             model.ErrTravelNotFound,
			},
			{
				name: "他人のTravelを更新しようとする",
				travelRepository: &model.Travel{
					ID:     travelID,
					UserID: 2,
				},
				wantErr: model.ErrUserNotAuthorized,
			},
			{
				name:                "TravelのRepositoryでエラーが発生する",
				travelRepositoryErr: databaseError,
				wantErr:             databaseError,
			},
			{
				name: "TravelPlanが存在しない",
				travelRepository: &model.Travel{
					ID:     travelID,
					UserID: userID,
				},
				travelPlanRepositoryErr: gorm.ErrRecordNotFound,
				wantErr:                 model.ErrTravelPlanNotFound,
			},
			{
				name: "TravelPlanのRepositoryでエラーが発生する",
				travelRepository: &model.Travel{
					ID:     travelID,
					UserID: userID,
				},
				travelPlanRepositoryErr: databaseError,
				wantErr:                 databaseError,
			},
			{
				name: "TravelPlanが指定されたTravelに属していない",
				travelRepository: &model.Travel{
					ID:     travelID,
					UserID: userID,
				},
				travelPlanRepository: &model.TravelPlan{
					ID:       travelPlanID,
					TravelID: 2,
				},
				wantErr: model.ErrTravelPlanNotBelongToTravel,
			},
			{
				name: "TravelPlanのUpdateでエラーが発生する",
				travelRepository: &model.Travel{
					ID:     travelID,
					UserID: userID,
				},
				travelPlanRepository: &model.TravelPlan{
					ID:       travelPlanID,
					TravelID: travelID,
				},
				updateErr: updateError,
				wantErr:   updateError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				repository := &mock.MockTravelPlanRepository{
					FindByIDFunc: func(id uint) (*model.TravelPlan, error) {
						return tt.travelPlanRepository, tt.travelPlanRepositoryErr
					},
					UpdateFunc: func(travelPlan *model.TravelPlan) error {
						return tt.updateErr
					},
				}

				travelRepository := &mock.MockTravelRepository{
					FindByIDFunc: func(id uint) (*model.Travel, error) {
						return tt.travelRepository, tt.travelRepositoryErr
					},
				}

				useCase := NewTravelPlanUsecase(travelRepository, repository)

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
