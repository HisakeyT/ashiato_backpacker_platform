package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

func TestHandler_GetTravel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name       string
		path       string
		travel     *model.Travel
		findError  error
		wantStatus int
		wantID     uint
		wantTitle  string
	}{
		{"own travel", "/travels/10", &model.Travel{ID: 10, UserID: 1, Title: "タイ旅行"}, nil, http.StatusOK, 10, "タイ旅行"},
		{"missing travel", "/travels/99", nil, model.ErrTravelNotFound, http.StatusNotFound, 0, ""},
		{"other user's travel", "/travels/10", &model.Travel{ID: 10, UserID: 2}, nil, http.StatusForbidden, 0, ""},
		{"invalid ID", "/travels/invalid", nil, nil, http.StatusBadRequest, 0, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &mock.MockTravelRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) { return test.travel, test.findError },
			}
			router := gin.New()
			handler := NewTravelHandler(usecase.NewTravelUsecase(repository))
			router.GET("/travels/:travel_id", func(c *gin.Context) {
				c.Set("userID", uint(1))
				handler.GetTravel(c)
			})
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, test.path, nil))
			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d: %s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			if test.wantStatus != http.StatusOK {
				var response struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				if response.Error == "" {
					t.Fatal("missing error in response")
				}
				return
			}
			var response TravelDetailResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if response.ID != test.wantID || response.Title != test.wantTitle {
				t.Fatalf("response = %+v", response)
			}
		})
	}
}
