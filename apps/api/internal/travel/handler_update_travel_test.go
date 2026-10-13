package travel

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestHandler_UpdateTravel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	databaseError := errors.New("database error")

	tests := []struct {
		name       string
		travelID   string
		request    string
		repository Repository
		wantStatus int
	}{
		{
			name:     "正常に更新できる",
			travelID: "1",
			request:  `{"title":"新しいタイトル"}`,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:        1,
						UserID:    1,
						Title:     "元のタイトル",
						StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC),
						IsPublic:  false,
					}, nil
				},
				UpdateFunc: func(travel *model.Travel) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "不正なJSON",
			travelID:   "1",
			request:    `{"title":}`,
			repository: &MockRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "travel_idが不正",
			travelID:   "abc",
			request:    `{"title":"新しいタイトル"}`,
			repository: &MockRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "開始日が不正",
			travelID: "1",
			request:  `{"start_date":"2026-99-99"}`,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:     1,
						UserID: 1,
					}, nil
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "終了日が不正",
			travelID: "1",
			request:  `{"end_date":"2026-99-99"}`,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:     1,
						UserID: 1,
					}, nil
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "Travelが存在しない",
			travelID: "1",
			request:  `{"title":"新しいタイトル"}`,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "他人のTravelを更新しようとする",
			travelID: "1",
			request:  `{"title":"新しいタイトル"}`,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:     1,
						UserID: 2,
					}, nil
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:     "開始日が終了日より後",
			travelID: "1",
			request:  `{"start_date":"2026-10-10"}`,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:        1,
						UserID:    1,
						StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC),
					}, nil
				},
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "Repositoryでエラーが発生する",
			travelID: "1",
			request:  `{"title":"新しいタイトル"}`,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return nil, databaseError
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			useCase := NewUseCase(tt.repository)
			handler := NewHandler(useCase)

			router.PATCH("/travels/:travel_id", func(c *gin.Context) {
				c.Set("userID", uint(1))
				handler.UpdateTravel(c)
			})

			req := httptest.NewRequest(
				http.MethodPatch,
				"/travels/"+tt.travelID,
				strings.NewReader(tt.request),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d",
					tt.wantStatus,
					rec.Code,
				)
			}
		})
	}
}
