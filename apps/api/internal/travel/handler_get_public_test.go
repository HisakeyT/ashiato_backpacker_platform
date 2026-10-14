package travel

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
	"github.com/gin-gonic/gin"
)

func TestHandler_GetPublicTravels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		repository repository.TravelRepository
		wantStatus int
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
			wantStatus: http.StatusOK,
		},
		{
			name: "公開Travelが0件",
			repository: &MockRepository{
				FindPublicTravelsFunc: func() ([]*model.PublicTravel, error) {
					return []*model.PublicTravel{}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "UseCaseでエラーが発生する",
			repository: &MockRepository{
				FindPublicTravelsFunc: func() ([]*model.PublicTravel, error) {
					return nil, errors.New("repository error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler := &Handler{
				useCase: &UseCase{
					repository: tt.repository,
				},
			}

			handler.GetPublicTravels(c)

			if w.Code != tt.wantStatus {
				t.Errorf("GetPublicTravels() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
