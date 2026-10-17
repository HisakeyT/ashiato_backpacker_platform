package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

func TestHandler_GetTravels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		repository repository.TravelRepository
		wantStatus int
		wantIDs    []uint
	}{
		{
			name: "自分の旅行を複数件取得できる",
			repository: &mock.MockTravelRepository{
				FindByUserIDFunc: func(userID uint) ([]*model.Travel, error) {
					return []*model.Travel{
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
					}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantIDs:    []uint{1, 2},
		},
		{
			name: "自分の旅行が0件でも正常に取得できる",
			repository: &mock.MockTravelRepository{
				FindByUserIDFunc: func(userID uint) ([]*model.Travel, error) {
					return []*model.Travel{}, nil
				},
			},
			wantStatus: http.StatusOK,
			wantIDs:    []uint{},
		},
		{
			name: "Repositoryでエラーが発生した場合",
			repository: &mock.MockTravelRepository{
				FindByUserIDFunc: func(userID uint) ([]*model.Travel, error) {
					return nil, errors.New("database error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			useCase := usecase.NewTravelUsecase(tt.repository)
			handler := NewTravelHandler(useCase)

			router.GET("/travels", func(c *gin.Context) {
				c.Set("userID", uint(1))
				handler.GetTravels(c)
			})

			req := httptest.NewRequest(
				http.MethodGet,
				"/travels",
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					rec.Code,
				)
			}
			if tt.wantStatus == http.StatusOK {
				var response []TravelResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				if len(response) != len(tt.wantIDs) {
					t.Fatalf("response count = %d, want %d", len(response), len(tt.wantIDs))
				}
				for index, travel := range response {
					if travel.ID != tt.wantIDs[index] || travel.Title == "" {
						t.Fatalf("unexpected travel response: %+v", travel)
					}
				}
			}
		})
	}
}
