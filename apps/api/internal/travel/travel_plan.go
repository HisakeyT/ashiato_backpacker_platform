package travel

import (
	"errors"
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
	Place     string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	SortOrder int       `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
