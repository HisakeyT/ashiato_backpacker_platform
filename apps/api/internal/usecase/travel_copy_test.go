package usecase

import (
	"errors"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
)

func TestUsecase_CopyTravel(t *testing.T) {
	errDB := errors.New("db error")

	newPublic := func() *model.Travel {
		return &model.Travel{
			ID: 10, UserID: 1, IsPublic: true, Title: "t",
			Plans: []model.TravelPlan{{ID: 100, TravelID: 10, Place: "p", SortOrder: 1}},
		}
	}

	tests := []struct {
		name       string
		userID     uint
		findResult *model.Travel
		findErr    error
		createErr  error
		wantErr    error
		wantSaved  bool
	}{
		{
			name:       "公開旅行を他人がコピーできる",
			userID:     2,
			findResult: newPublic(),
			wantSaved:  true},
		{
			name:       "自分の非公開旅行はコピーできる",
			userID:     1,
			findResult: &model.Travel{ID: 10, UserID: 1, IsPublic: false, Title: "mine"},
			wantSaved:  true,
		},
		{
			name:       "他人の非公開旅行はNotFound",
			userID:     2,
			findResult: &model.Travel{ID: 10, UserID: 1, IsPublic: false},
			wantErr:    model.ErrTravelNotFound,
		},
		{
			name:    "存在しない旅行はNotFound",
			userID:  2,
			findErr: model.ErrTravelNotFound,
			wantErr: model.ErrTravelNotFound,
		},
		{
			name:    "取得時のDBエラーはそのまま返す",
			userID:  2,
			findErr: errDB,
			wantErr: errDB,
		},
		{
			name:       "保存時のDBエラーはそのまま返す",
			userID:     2,
			findResult: newPublic(),
			createErr:  errDB,
			wantErr:    errDB,
			wantSaved:  true,
		},
		{
			name:       "プランが無い旅行もコピーできる",
			userID:     2,
			findResult: &model.Travel{ID: 10, UserID: 1, IsPublic: true, Title: "t"},
			wantSaved:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var saved *model.Travel
			repo := &mock.MockTravelRepository{
				FindByIDWithPlansFunc: func(id uint) (*model.Travel, error) {
					return tt.findResult, tt.findErr
				},
				CreateWithPlansFunc: func(travel *model.Travel) error {
					saved = travel
					if tt.createErr != nil {
						return tt.createErr
					}
					travel.ID = 20 // DB の採番を模擬
					return nil
				},
			}
			uc := NewTravelUsecase(repo)

			got, err := uc.CopyTravel(tt.userID, 10)

			if (saved != nil) != tt.wantSaved {
				t.Fatalf("saved = %v, wantSaved %v", saved != nil, tt.wantSaved)
			}
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				if got != nil {
					t.Errorf("got = %+v, want nil on error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.ID != 20 {
				t.Errorf("ID = %d, want 20", got.ID)
			}
			if got.UserID != tt.userID {
				t.Errorf("UserID = %d, want %d", got.UserID, tt.userID)
			}
			if got.CopiedFromTravelID == nil || *got.CopiedFromTravelID != 10 {
				t.Errorf("CopiedFromTravelID = %v, want 10", got.CopiedFromTravelID)
			}
			if got.IsPublic {
				t.Error("copied travel must start as private")
			}
			if len(got.Plans) != len(tt.findResult.Plans) {
				t.Errorf("len(Plans) = %d, want %d", len(got.Plans), len(tt.findResult.Plans))
			}
		})
	}
}
