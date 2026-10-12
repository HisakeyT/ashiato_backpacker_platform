package travel

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestHandler_UpdateTravelPlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		travelID   string
		planID     string
		request    string
		repository PlanRepository
		wantStatus int
	}{
		{
			name:     "日付を更新できる",
			travelID: "1",
			planID:   "1",
			request:  `{"start_date":"2026-10-02"}`,
			repository: &MockPlanRepository{
				FindByIDFunc: func(id uint) (*TravelPlan, error) {
					return &TravelPlan{
						ID:        1,
						TravelID:  1,
						StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
						Place:     "バンコク",
						Content:   "ワット・ポー",
						SortOrder: 1,
					}, nil
				},
				UpdateFunc: func(travelPlan *TravelPlan) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "Placeを更新できる",
			travelID: "1",
			planID:   "1",
			request:  `{"place":"チェンマイ"}`,
			repository: &MockPlanRepository{
				FindByIDFunc: func(id uint) (*TravelPlan, error) {
					return &TravelPlan{
						ID:        1,
						TravelID:  1,
						StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
						Place:     "バンコク",
						Content:   "ワット・ポー",
						SortOrder: 1,
					}, nil
				},
				UpdateFunc: func(travelPlan *TravelPlan) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "Contentを更新できる",
			travelID: "1",
			planID:   "1",
			request:  `{"content":"ワット・プラ・シン"}`,
			repository: &MockPlanRepository{
				FindByIDFunc: func(id uint) (*TravelPlan, error) {
					return &TravelPlan{
						ID:        1,
						TravelID:  1,
						StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
						Place:     "バンコク",
						Content:   "ワット・ポー",
						SortOrder: 1,
					}, nil
				},
				UpdateFunc: func(travelPlan *TravelPlan) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "SortOrderを更新できる",
			travelID: "1",
			planID:   "1",
			request:  `{"sort_order":2}`,
			repository: &MockPlanRepository{
				FindByIDFunc: func(id uint) (*TravelPlan, error) {
					return &TravelPlan{
						ID:        1,
						TravelID:  1,
						StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
						Place:     "バンコク",
						Content:   "ワット・ポー",
						SortOrder: 1,
					}, nil
				},
				UpdateFunc: func(travelPlan *TravelPlan) error {
					return nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "不正なJSON",
			travelID:   "1",
			planID:     "1",
			request:    `{"date":}`,
			repository: &MockPlanRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "travel_idが不正",
			travelID:   "abc",
			planID:     "1",
			request:    `{}`,
			repository: &MockPlanRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "travel_plan_idが不正",
			travelID:   "1",
			planID:     "abc",
			request:    `{}`,
			repository: &MockPlanRepository{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "日付が不正",
			travelID:   "1",
			planID:     "1",
			request:    `{"start_date":"2026-99-99"}`,
			repository: &MockPlanRepository{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			// TravelRepositoryはUseCaseで必要なので、
			// 実際のプロジェクトのMockTravelRepositoryに合わせる
			travelRepository := &MockTravelRepository{
				FindByIDFunc: func(id uint) (*Travel, error) {
					return &Travel{
						ID:     1,
						UserID: 1,
					}, nil
				},
			}

			useCase := NewPlanUseCase(travelRepository, tt.repository)
			handler := NewPlanHandler(useCase)

			router.PATCH(
				"/travels/:travel_id/plans/:travel_plan_id",
				func(c *gin.Context) {
					c.Set("userID", uint(1))
					handler.UpdateTravelPlan(c)
				},
			)

			req := httptest.NewRequest(
				http.MethodPatch,
				"/travels/"+tt.travelID+"/plans/"+tt.planID,
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
