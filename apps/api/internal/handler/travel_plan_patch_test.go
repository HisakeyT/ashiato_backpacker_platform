package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

func TestHandler_UpdateTravelPlan_PatchFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	startDate := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 0, 1)
	country := "Japan"
	state := "Tokyo"
	city := "Tokyo"
	latitude := 35.0
	longitude := 139.0

	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantStart  bool
		wantEnd    bool
	}{
		{"null clears optional fields", `{"end_date":null,"state":null}`, http.StatusOK, true, false},
		{"omitted date fields are preserved", `{}`, http.StatusOK, true, true},
		{"null clears both dates in planning", `{"start_date":null,"end_date":null}`, http.StatusOK, false, false},
		{"unpaired coordinate is rejected", `{"latitude":35}`, http.StatusBadRequest, true, false},
		{"location values replace prior values", `{"city":"Kyoto"}`, http.StatusOK, true, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var saved *model.TravelPlan
			planRepository := &mock.MockTravelPlanRepository{
				FindByIDFunc: func(id uint) (*model.TravelPlan, error) {
					return &model.TravelPlan{
						ID: 2, TravelID: 1, StartDate: &startDate, EndDate: &endDate,
						Place: "place", Content: "content",
						Location: model.Location{Country: &country, State: &state, City: &city, Latitude: &latitude, Longitude: &longitude},
					}, nil
				},
				UpdateFunc: func(plan *model.TravelPlan) error { saved = plan; return nil },
			}
			travelRepository := &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{ID: 1, UserID: 1, StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)}, nil
				},
			}
			router := gin.New()
			handler := NewTravelPlanHandler(usecase.NewTravelPlanUsecase(travelRepository, planRepository))
			router.PATCH("/travels/:travel_id/plans/:travel_plan_id", func(c *gin.Context) {
				c.Set("userID", uint(1))
				handler.UpdateTravelPlan(c)
			})
			request := httptest.NewRequest(http.MethodPatch, "/travels/1/plans/2", strings.NewReader(test.body))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			if test.wantStatus != http.StatusOK {
				if saved != nil {
					t.Fatal("invalid patch was saved")
				}
				return
			}
			if saved == nil {
				t.Fatal("plan was not saved")
			}
			if (saved.StartDate != nil) != test.wantStart ||
				(saved.EndDate != nil) != test.wantEnd || saved.Location.Country != nil ||
				saved.Location.State != nil || saved.Location.Latitude != nil || saved.Location.Longitude != nil {
				t.Fatalf("unexpected saved plan: %+v", saved)
			}
			if test.body == `{"city":"Kyoto"}` {
				if saved.Location.City == nil || *saved.Location.City != "Kyoto" {
					t.Fatalf("unexpected city: %+v", saved.Location)
				}
			} else if saved.Location.City != nil {
				t.Fatalf("unexpected city: %+v", saved.Location)
			}
			var response TravelPlanResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(response.City, saved.Location.City) ||
				(response.StartDate != nil) != test.wantStart || (response.EndDate != nil) != test.wantEnd {
				t.Fatalf("unexpected response: %+v", response)
			}
		})
	}
}
