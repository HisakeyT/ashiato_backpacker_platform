package travel

import (
	"errors"
	"gorm.io/gorm"
	"reflect"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
)

func TestUseCase_GetTravel(t *testing.T) {
	userID := uint(1)
	sourceID := uint(10)

	tests := []struct {
		name                string
		travelID            uint
		repository          *model.Travel
		repositoryError     error
		copySource          *model.CopySource
		copySourceError     error
		wantCopySourceCalls int
		want                *model.TravelDetail
		wantErr             error
	}{
		{
			name:     "自分の旅行を取得できる",
			travelID: 1,
			repository: &model.Travel{
				ID:     1,
				UserID: userID,
				Title:  "旅行1",
			},
			want: &model.TravelDetail{
				Travel: &model.Travel{ID: 1, UserID: userID, Title: "旅行1"},
			},
		},
		{
			name:     "コピー元があれば CopiedFrom に入る",
			travelID: 3,
			repository: &model.Travel{
				ID:                 3,
				UserID:             userID,
				Title:              "コピーした旅行",
				CopiedFromTravelID: &sourceID,
			},
			copySource:          &model.CopySource{TravelID: sourceID, AuthorNickname: "たろう"},
			wantCopySourceCalls: 1,
			want: &model.TravelDetail{
				Travel: &model.Travel{
					ID:                 3,
					UserID:             userID,
					Title:              "コピーした旅行",
					CopiedFromTravelID: &sourceID,
				},
				CopiedFrom: &model.CopySource{TravelID: sourceID, AuthorNickname: "たろう"},
			},
		},
		{
			name:     "コピー元が非公開なら CopiedFrom は nil",
			travelID: 3,
			repository: &model.Travel{
				ID:                 3,
				UserID:             userID,
				Title:              "コピーした旅行",
				CopiedFromTravelID: &sourceID,
			},
			copySource:          nil,
			wantCopySourceCalls: 1,
			want: &model.TravelDetail{
				Travel: &model.Travel{
					ID:                 3,
					UserID:             userID,
					Title:              "コピーした旅行",
					CopiedFromTravelID: &sourceID,
				},
			},
		},
		{
			name:            "旅行が存在しない",
			travelID:        999,
			repositoryError: gorm.ErrRecordNotFound,
			wantErr:         model.ErrTravelNotFound,
		},
		{
			name:     "他人の旅行を取得しようとした",
			travelID: 2,
			repository: &model.Travel{
				ID:     2,
				UserID: 2,
				Title:  "他人の旅行",
			},
			wantErr: model.ErrUserNotAuthorized,
		},
		{
			name:            "Repositoryでエラーが発生した場合、UseCaseでもエラーを返す",
			travelID:        1,
			repositoryError: errors.New("database error"),
			wantErr:         errors.New("database error"),
		},
		{
			name:     "FindCopySourceでエラーが発生した場合、エラーを返す",
			travelID: 3,
			repository: &model.Travel{
				ID:                 3,
				UserID:             userID,
				CopiedFromTravelID: &sourceID,
			},
			copySourceError:     errors.New("copy source error"),
			wantCopySourceCalls: 1,
			wantErr:             errors.New("copy source error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copySourceCalls := 0
			repository := &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return tt.repository, tt.repositoryError
				},
				FindCopySourceFunc: func(id uint) (*model.CopySource, error) {
					copySourceCalls++
					if id != sourceID {
						t.Errorf("FindCopySource called with %d, want %d", id, sourceID)
					}
					return tt.copySource, tt.copySourceError
				},
			}

			useCase := NewUseCase(repository)
			got, err := useCase.GetTravel(userID, tt.travelID)

			if copySourceCalls != tt.wantCopySourceCalls {
				t.Fatalf("FindCopySource calls = %d, want %d", copySourceCalls, tt.wantCopySourceCalls)
			}

			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}
