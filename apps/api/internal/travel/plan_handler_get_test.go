package travel

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestHandler_GetTravelPlans(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		userID          uint
		travelID        string
		setupRepository func() (*MockTravelRepository, *MockPlanRepository)
		wantStatus      int
	}{
		{
			name:     "正常に取得できる",
			userID:   1,
			travelID: "10",
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*model.Travel, error) {
						return &model.Travel{
							ID:     10,
							UserID: 1,
						}, nil
					},
				}

				travelPlanRepository := &MockPlanRepository{
					FindByTravelIDFunc: func(id uint) ([]model.TravelPlan, error) {
						return []model.TravelPlan{
							{
								ID:        1,
								TravelID:  10,
								Place:     "バンコク",
								Content:   "ワット・ポーを観光",
								SortOrder: 1,
							},
							{
								ID:        2,
								TravelID:  10,
								Place:     "バンコク",
								Content:   "王宮を観光",
								SortOrder: 2,
							},
						}, nil
					},
				}

				return travelRepository, travelPlanRepository
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "TravelPlanが0件でも取得できる",
			userID:   1,
			travelID: "10",
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*model.Travel, error) {
						return &model.Travel{
							ID:     10,
							UserID: 1,
						}, nil
					},
				}

				travelPlanRepository := &MockPlanRepository{
					FindByTravelIDFunc: func(id uint) ([]model.TravelPlan, error) {
						return []model.TravelPlan{}, nil
					},
				}

				return travelRepository, travelPlanRepository
			},
			wantStatus: http.StatusOK,
		},
		{
			name:     "旅行が存在しない",
			userID:   1,
			travelID: "999",
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*model.Travel, error) {
						return nil, gorm.ErrRecordNotFound
					},
				}

				travelPlanRepository := &MockPlanRepository{}

				return travelRepository, travelPlanRepository
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "他人の旅行は取得できない",
			userID:   1,
			travelID: "10",
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				travelRepository := &MockTravelRepository{
					FindByIDFunc: func(id uint) (*model.Travel, error) {
						return &model.Travel{
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
			name:     "travel_idが不正",
			userID:   1,
			travelID: "abc",
			setupRepository: func() (*MockTravelRepository, *MockPlanRepository) {
				return &MockTravelRepository{}, &MockPlanRepository{}
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
				http.MethodGet,
				"/travels/"+tt.travelID+"/plans",
				nil,
			)

			handler.GetTravelPlans(c)

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
