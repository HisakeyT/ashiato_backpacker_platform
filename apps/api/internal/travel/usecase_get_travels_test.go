package travel

import (
	"errors"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
)

func TestUseCase_GetTravels(t *testing.T) {
	userID := uint(1)

	tests := []struct {
		name            string
		repository      []*model.Travel
		repositoryError error
		wantLen         int
	}{
		{
			name: "自分の旅行を複数件取得できる",
			repository: []*model.Travel{
				{
					ID:     1,
					UserID: userID,
					Title:  "旅行1",
				},
				{
					ID:     2,
					UserID: userID,
					Title:  "旅行2",
				},
			},
			wantLen: 2,
		},
		{
			name:       "自分の旅行が0件でも正常に値を返す",
			repository: []*model.Travel{},
			wantLen:    0,
		},
		{
			name:            "RepositoryのFindByUserIDでエラーが発生した場合、UseCaseでもエラーを返す",
			repositoryError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &MockRepository{
				FindByUserIDFunc: func(userID uint) ([]*model.Travel, error) {
					return tt.repository, tt.repositoryError
				},
			}

			useCase := NewUseCase(repository)
			travels, err := useCase.GetTravels(userID)

			if tt.repositoryError != nil {
				if !errors.Is(err, tt.repositoryError) {
					t.Fatalf("expected error %v, got %v", tt.repositoryError, err)
				}

				return
			}

			if len(travels) != tt.wantLen {
				t.Fatalf("expected %d travels, got %d", tt.wantLen, len(travels))
			}
		})
	}
}
