package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

func TestTravelPlanOutsideTravelDatesResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	date := func(day int) time.Time { return time.Date(2026, 10, day, 0, 0, 0, 0, time.UTC) }
	for _, tt := range []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{"create before travel", http.MethodPost, "/travels/1/plans", `{"start_date":"2026-09-30","place":"Kyoto","content":"Visit"}`},
		{"create end after travel", http.MethodPost, "/travels/1/plans", `{"start_date":"2026-10-02","end_date":"2026-10-04","place":"Kyoto","content":"Visit"}`},
		{"update before travel", http.MethodPatch, "/travels/1/plans/2", `{"start_date":"2026-09-30"}`},
		{"update end after travel", http.MethodPatch, "/travels/1/plans/2", `{"end_date":"2026-10-04"}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			travelRepository := &mock.MockTravelRepository{FindByIDFunc: func(uint) (*model.Travel, error) {
				return &model.Travel{ID: 1, UserID: 1, StartDate: date(1), EndDate: date(3)}, nil
			}}
			planRepository := &mock.MockTravelPlanRepository{FindByIDFunc: func(uint) (*model.TravelPlan, error) {
				return &model.TravelPlan{ID: 2, TravelID: 1, StartDate: date(2), Place: "Kyoto", Content: "Visit"}, nil
			}}
			handler := NewTravelPlanHandler(usecase.NewTravelPlanUsecase(travelRepository, planRepository))
			router := gin.New()
			router.POST("/travels/:travel_id/plans", func(c *gin.Context) { c.Set("userID", uint(1)); handler.CreateTravelPlan(c) })
			router.PATCH("/travels/:travel_id/plans/:travel_plan_id", func(c *gin.Context) { c.Set("userID", uint(1)); handler.UpdateTravelPlan(c) })
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400: %s", recorder.Code, recorder.Body.String())
			}
			var response struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if response.Error != model.ErrTravelPlanOutsideTravelDates.Error() {
				t.Fatalf("error = %q, want %q", response.Error, model.ErrTravelPlanOutsideTravelDates)
			}
		})
	}
}
