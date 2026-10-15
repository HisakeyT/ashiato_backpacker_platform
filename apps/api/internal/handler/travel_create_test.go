package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

func TestHandler_CreateTravel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		request    string
		repository repository.TravelRepository
		wantStatus int
	}{
		{
			name:    "正常に作成できる",
			request: `{"title":"東南アジア3週間","start_date":"2026-10-01","end_date":"2026-10-21","is_public":true}`,
			repository: &mock.MockTravelRepository{
				CreateFunc: func(travel *model.Travel) error {
					travel.ID = 1
					return nil
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "不正なJSON",
			request:    `{"title":}`,
			repository: &mock.MockTravelRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "開始日が不正",
			request:    `{"title":"東南アジア3週間","start_date":"2026-99-99","end_date":"2026-10-21","is_public":true}`,
			repository: &mock.MockTravelRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "終了日が不正",
			request:    `{"title":"東南アジア3週間","start_date":"2026-10-01","end_date":"2026-99-99","is_public":true}`,
			repository: &mock.MockTravelRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "開始日が終了日より後",
			request:    `{"title":"東南アジア3週間","start_date":"2026-10-21","end_date":"2026-10-01","is_public":true}`,
			repository: &mock.MockTravelRepository{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			useCase := usecase.NewTravelUsecase(tt.repository)
			handler := NewTravelHandler(useCase)

			router.POST("/travels", func(c *gin.Context) {
				c.Set("userID", uint(1))
				handler.CreateTravel(c)
			})

			req := httptest.NewRequest(
				http.MethodPost,
				"/travels",
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
