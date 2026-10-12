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
	EndDate   *time.Time
	Place     string `gorm:"not null"`
	Country   *string
	State     *string
	City      *string
	Latitude  *float64
	Longitude *float64
	Content   string `gorm:"not null"`
	SortOrder int    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
