package usecase

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
	"gorm.io/gorm"
)

func TestUsecase_UpdateTravel(t *testing.T) {
	userID := uint(1)
	travelID := uint(1)

	startDate := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)

	newTitle := "新しいタイトル"
	newStartDate := time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC)
	newEndDate := time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC)
	newIsPublic := true

	tests := []struct {
		name        string
		updateInput UpdateTravelInput
		want        *model.Travel
	}{
		{
			name: "タイトルを更新する",
			updateInput: UpdateTravelInput{
				Title: &newTitle,
			},
			want: &model.Travel{
				ID:        travelID,
				UserID:    userID,
				Title:     newTitle,
				StartDate: startDate,
				EndDate:   endDate,
				IsPublic:  false,
			},
		},
		{
			name: "開始日を更新する",
			updateInput: UpdateTravelInput{
				StartDate: &newStartDate,
			},
			want: &model.Travel{
				ID:        travelID,
				UserID:    userID,
				Title:     "元のタイトル",
				StartDate: newStartDate,
				EndDate:   endDate,
				IsPublic:  false,
			},
		},
		{
			name: "終了日を更新する",
			updateInput: UpdateTravelInput{
				EndDate: &newEndDate,
			},
			want: &model.Travel{
				ID:        travelID,
				UserID:    userID,
				Title:     "元のタイトル",
				StartDate: startDate,
				EndDate:   newEndDate,
				IsPublic:  false,
			},
		},
		{
			name: "公開状態を更新する",
			updateInput: UpdateTravelInput{
				IsPublic: &newIsPublic,
			},
			want: &model.Travel{
				ID:        travelID,
				UserID:    userID,
				Title:     "元のタイトル",
				StartDate: startDate,
				EndDate:   endDate,
				IsPublic:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:        travelID,
						UserID:    userID,
						Title:     "元のタイトル",
						StartDate: startDate,
						EndDate:   endDate,
						IsPublic:  false,
					}, nil
				},
				UpdateFunc: func(travel *model.Travel) error {
					return nil
				},
			}

			useCase := NewTravelUsecase(repository)

			travel, err := useCase.UpdateTravel(
				userID,
				travelID,
				tt.updateInput,
			)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(travel, tt.want) {
				t.Fatalf("expected travel %+v, got %+v", tt.want, travel)
			}
		})
	}
}

func TestUsecase_UpdateTravel_Error(t *testing.T) {
	userID := uint(1)
	travelID := uint(1)
	databaseError := errors.New("database error")
	updateError := errors.New("update error")

	invalidStartDate := time.Date(
		2026, 10, 10, 0, 0, 0, 0, time.UTC,
	)
	startDate := time.Date(
		2026, 10, 1, 0, 0, 0, 0, time.UTC,
	)
	endDate := time.Date(
		2026, 10, 5, 0, 0, 0, 0, time.UTC,
	)

	tests := []struct {
		name          string
		updateInput   UpdateTravelInput
		repository    *model.Travel
		repositoryErr error
		updateErr     error
		wantErr       error
	}{
		{
			name:          "Travelが存在しない",
			repositoryErr: gorm.ErrRecordNotFound,
			wantErr:       model.ErrTravelNotFound,
		},
		{
			name: "他人のTravelを更新しようとする",
			repository: &model.Travel{
				ID:     travelID,
				UserID: 2,
				Title:  "他人の旅行",
			},
			wantErr: model.ErrUserNotAuthorized,
		},
		{
			name:          "RepositoryのFindByIDでエラーが発生する",
			repositoryErr: databaseError,
			wantErr:       databaseError,
		},
		{
			name:      "RepositoryのUpdateでエラーが発生する",
			updateErr: updateError,
			wantErr:   updateError,
		},
		{
			name: "開始日が終了日より後になる",
			updateInput: UpdateTravelInput{
				StartDate: &invalidStartDate,
			},
			repository: &model.Travel{
				ID:        travelID,
				UserID:    userID,
				Title:     "元のタイトル",
				StartDate: startDate,
				EndDate:   endDate,
			},
			wantErr: ErrInvalidDateRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					if tt.repositoryErr != nil {
						return nil, tt.repositoryErr
					}

					if tt.repository != nil {
						return tt.repository, nil
					}

					return &model.Travel{
						ID:     travelID,
						UserID: userID,
						Title:  "元のタイトル",
					}, nil
				},
				UpdateFunc: func(travel *model.Travel) error {
					return tt.updateErr
				},
			}

			useCase := NewTravelUsecase(repository)

			_, err := useCase.UpdateTravel(
				userID,
				travelID,
				tt.updateInput,
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
}
