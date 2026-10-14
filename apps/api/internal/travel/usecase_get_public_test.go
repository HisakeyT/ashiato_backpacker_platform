package travel

import (
	"errors"
	"reflect"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
)

func TestUseCase_GetPublicTravels(t *testing.T) {
	tests := []struct {
		name        string
		repository  repository.Repository
		wantTravels []*model.PublicTravel
		wantErr     bool
	}{
		{
			name: "公開Travelを複数取得できる",
			repository: &MockRepository{
				FindPublicTravelsFunc: func() ([]*model.PublicTravel, error) {
					return []*model.PublicTravel{
						{
							ID:       1,
							UserID:   1,
							Title:    "タイ旅行",
							IsPublic: true,
						},
						{
							ID:       2,
							UserID:   2,
							Title:    "ベトナム旅行",
							IsPublic: true,
						},
					}, nil
				},
			},
			wantTravels: []*model.PublicTravel{
				{
					ID:       1,
					UserID:   1,
					Title:    "タイ旅行",
					IsPublic: true,
				},
				{
					ID:       2,
					UserID:   2,
					Title:    "ベトナム旅行",
					IsPublic: true,
				},
			},
		},
		{
			name: "公開Travelが0件",
			repository: &MockRepository{
				FindPublicTravelsFunc: func() ([]*model.PublicTravel, error) {
					return []*model.PublicTravel{}, nil
				},
			},
			wantTravels: []*model.PublicTravel{},
		},
		{
			name: "Repositoryエラー",
			repository: &MockRepository{
				FindPublicTravelsFunc: func() ([]*model.PublicTravel, error) {
					return nil, errors.New("repository error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &UseCase{
				repository: tt.repository,
			}

			got, err := uc.GetPublicTravels()

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetPublicTravels() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.wantTravels) {
				t.Errorf("GetPublicTravels() = %v, want %v", got, tt.wantTravels)
			}
		})
	}
}
