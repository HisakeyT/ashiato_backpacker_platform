package travel

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestHandler_CreateTravelPlan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		userID          uint
		travelID        string
		body            string
		setupRepository func() (*MockTravelRepository, *MockPlanRepository)
		wantStatus      int
	}{
		{
			name:     "正常に作成できる",
			userID:   1,
			travelID: "10",
			body: `{
				"start_date": "2026-10-01",
				"place": "バンコク",
				"content": "ワット・ポーを観光",
				"sort_order": 1
			}`,
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*Travel, error) {
						return &Travel{
							ID:     10,
							UserID: 1,
						}, nil
					},
				}

				travelPlanRepository := &MockPlanRepository{
					CreateFunc: func(travelPlan *TravelPlan) error {
						return nil
					},
				}

				return travelRepository, travelPlanRepository
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:     "旅行が存在しない",
			userID:   1,
			travelID: "999",
			body: `{
				"start_date": "2026-10-01",
				"place": "バンコク",
				"content": "ワット・ポーを観光",
				"sort_order": 1
			}`,
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*Travel, error) {
						return nil, gorm.ErrRecordNotFound
					},
				}

				travelPlanRepository := &MockPlanRepository{}

				return travelRepository, travelPlanRepository
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "他人の旅行には追加できない",
			userID:   1,
			travelID: "10",
			body: `{
				"start_date": "2026-10-01",
				"place": "バンコク",
				"content": "ワット・ポーを観光",
				"sort_order": 1
			}`,
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*Travel, error) {
						return &Travel{
							ID:     10,
							UserID: 999,
						}, nil
					},
				}

				travelPlanRepository := &MockPlanRepository{}

				return travelRepository, travelPlanRepository
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name:     "不正なJSON",
			userID:   1,
			travelID: "10",
			body:     `{"date":`,
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				return &MockTravelRepository{}, &MockPlanRepository{}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "travel_idが不正",
			userID:   1,
			travelID: "abc",
			body: `{
			"start_date": "2026-10-01",
				"place": "バンコク",
				"content": "ワット・ポーを観光",
				"sort_order": 1
			}`,
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				return &MockTravelRepository{}, &MockPlanRepository{}
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:     "start_dateが不正",
			userID:   1,
			travelID: "10",
			body: `{
				"start_date": "2026/10/01",
				"place": "バンコク",
				"content": "ワット・ポーを観光",
				"sort_order": 1
			}`,
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				return &MockTravelRepository{
					FindByIDFunc: func(id uint) (*Travel, error) {
						return &Travel{
							ID:     10,
							UserID: 1,
						}, nil
					},
				}, &MockPlanRepository{}
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			travelRepository, travelPlanRepository := tt.setupRepository()

			useCase := NewPlanUseCase(
				travelRepository,
				travelPlanRepository,
			)

			handler := NewPlanHandler(useCase)

			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)

			c.Params = gin.Params{
				{
					Key:   "travel_id",
					Value: tt.travelID,
				},
			}
			c.Set("userID", tt.userID)

			c.Request = httptest.NewRequest(
				http.MethodPost,
				"/travels/"+tt.travelID+"/plans",
				strings.NewReader(tt.body),
			)
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateTravelPlan(c)

			if recorder.Code != tt.wantStatus {
				t.Errorf(
					"status = %d, want %d",
					recorder.Code,
					tt.wantStatus,
				)
			}
		})
	}
}
