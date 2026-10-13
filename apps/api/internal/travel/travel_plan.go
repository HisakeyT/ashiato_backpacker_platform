package travel

import (
	"errors"
	"github.com/HisakeyT/backpacker-platform/internal/model"
	"time"
)

var (
	ErrTravelPlanNotFound          = errors.New("travel plan not found")
	ErrTravelPlanNotBelongToTravel = errors.New("travel plan does not belong to the specified travel")
)

type TravelPlan struct {
	ID        uint      `gorm:"primaryKey"`
	TravelID  uint      `gorm:"not null;index"`
	StartDate time.Time `gorm:"not null"`
	EndDate   *time.Time
	Place     string         `gorm:"not null"`
	Location  model.Location `gorm:"embedded"`
	Content   string         `gorm:"not null"`
	SortOrder int            `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
