package travel

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
	"github.com/gin-gonic/gin"
)

func TestHandler_DeleteTravel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		travelID   string
		userID     uint
		repository repository.Repository
		wantStatus int
	}{
		{
			name:     "正常に削除できる",
			travelID: "1",
			userID:   1,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:     1,
						UserID: 1,
					}, nil
				},
				DeleteFunc: func(id uint) error {
					return nil
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:     "Travelが存在しない",
			travelID: "999",
			userID:   1,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return nil, gorm.ErrRecordNotFound
				},
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:     "他ユーザーのTravel",
			travelID: "1",
			userID:   1,
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
			name:     "FindByIDでRepositoryエラー",
			travelID: "1",
			userID:   1,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return nil, errors.New("database error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:     "DeleteでRepositoryエラー",
			travelID: "1",
			userID:   1,
			repository: &MockRepository{
				FindByIDFunc: func(id uint) (*model.Travel, error) {
					return &model.Travel{
						ID:     1,
						UserID: 1,
					}, nil
				},
				DeleteFunc: func(id uint) error {
					return errors.New("delete error")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "travel_idが不正",
			travelID:   "abc",
			userID:     1,
			repository: &MockRepository{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			useCase := NewUseCase(tt.repository)
			handler := NewHandler(useCase)

			router.DELETE("/travels/:travel_id", func(c *gin.Context) {
				c.Set("userID", tt.userID)
				handler.DeleteTravel(c)
			})

			req := httptest.NewRequest(
				http.MethodDelete,
				"/travels/"+tt.travelID,
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
		})
	}
}
