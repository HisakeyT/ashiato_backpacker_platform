package model

import (
	"errors"
	"time"
)

var (
	ErrTravelPlanNotFound           = errors.New("travel plan not found")
	ErrTravelPlanNotBelongToTravel  = errors.New("travel plan does not belong to the specified travel")
	ErrInvalidTravelPlanText        = errors.New("place and content must not be empty")
	ErrTravelPlanOutsideTravelDates = errors.New("travel plan dates must be within travel dates")
)

type TravelPlan struct {
	ID        uint      `gorm:"primaryKey"`
	TravelID  uint      `gorm:"not null;index"`
	StartDate time.Time `gorm:"not null"`
	EndDate   *time.Time
	Place     string   `gorm:"not null"`
	Location  Location `gorm:"embedded"`
	Content   string   `gorm:"not null"`
	SortOrder int      `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
