package travel

import (
	"errors"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"gorm.io/gorm"
)

func TestUseCase_CopyTravel(t *testing.T) {
	tests := []struct {
		name       string
		userID     uint
		travelID   uint
		repository *MockRepository
		wantErr    error
		wantID     uint
	}{
		{
			name:     "公開旅行を他人がコピーできる",
			userID:   2,
			travelID: 10,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{ID: 10, UserID: 1, IsPublic: true}, nil
				},
				CopyTravelFunc: func(src *model.Travel, newOwnerID uint) (*model.Travel, error) {
					return &model.Travel{ID: 20, UserID: newOwnerID, CopiedFromTravelID: &src.ID}, nil
				},
			},
			wantID: 20,
		},
		{
			name:     "自分の非公開旅行はコピーできる",
			userID:   1,
			travelID: 10,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{ID: 10, UserID: 1, IsPublic: false}, nil
				},
				CopyTravelFunc: func(src *model.Travel, newOwnerID uint) (*model.Travel, error) {
					return &model.Travel{ID: 21, UserID: newOwnerID}, nil
				},
			},
			wantID: 21,
		},
		{
			name:     "他人の非公開旅行はNotFound",
			userID:   2,
			travelID: 10,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{ID: 10, UserID: 1, IsPublic: false}, nil
				},
			},
			wantErr: model.ErrTravelNotFound,
		},
		{
			name:     "存在しない旅行はNotFound",
			userID:   2,
			travelID: 99,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			wantErr: model.ErrTravelNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewUseCase(tt.repository)

			got, err := uc.CopyTravel(tt.userID, tt.travelID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", got.ID, tt.wantID)
			}
		})
	}
}
