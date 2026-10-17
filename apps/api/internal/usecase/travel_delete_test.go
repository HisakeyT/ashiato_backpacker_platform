package usecase

import (
	"errors"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
	mock "github.com/HisakeyT/backpacker-platform/internal/repository/mock"
)

func TestUsecase_DeleteTravel(t *testing.T) {
	tests := []struct {
		name       string
		userID     uint
		travelID   uint
		repository repository.TravelRepository
		wantErr    error
	}{
		{
			name:     "自分のTravelを削除できる",
			userID:   1,
			travelID: 10,
			repository: &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:     10,
						UserID: 1,
					}, nil
				},
				DeleteFunc: func(id uint) error {
					return nil
				},
			},
			wantErr: nil,
		},
		{
			name:     "Travelが存在しない",
			userID:   1,
			travelID: 10,
			repository: &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return nil, model.ErrTravelNotFound
				},
				DeleteFunc: func(id uint) error {
					return nil
				},
			},
			wantErr: model.ErrTravelNotFound,
		},
		{
			name:     "他ユーザーのTravel",
			userID:   1,
			travelID: 10,
			repository: &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:     10,
						UserID: 2,
					}, nil
				},
				DeleteFunc: func(id uint) error {
					return nil
				},
			},
			wantErr: model.ErrUserNotAuthorized,
		},
		{
			name:     "FindByIDでRepositoryエラー",
			userID:   1,
			travelID: 10,
			repository: &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return nil, errors.New("database error")
				},
				DeleteFunc: func(id uint) error {
					return nil
				},
			},
			wantErr: errors.New("database error"),
		},
		{
			name:     "DeleteでRepositoryエラー",
			userID:   1,
			travelID: 10,
			repository: &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:     10,
						UserID: 1,
					}, nil
				},
				DeleteFunc: func(id uint) error {
					return errors.New("delete error")
				},
			},
			wantErr: errors.New("delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &TravelUsecase{
				repository: tt.repository,
			}

			err := u.DeleteTravel(tt.userID, tt.travelID)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("DeleteTravel() error = %v, wantErr nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("DeleteTravel() error = nil, wantErr %v", tt.wantErr)
			}

			if err.Error() != tt.wantErr.Error() {
				t.Errorf("DeleteTravel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
