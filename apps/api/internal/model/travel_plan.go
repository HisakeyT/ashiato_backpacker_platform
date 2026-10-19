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
	ErrTravelPlanStartDateRequired  = errors.New("travel plan start date is required")
	ErrTravelPlanEndWithoutStart    = errors.New("travel plan end date requires start date")
	ErrInvalidTravelPlanOrder       = errors.New("plan order must contain every plan in this travel exactly once")
	ErrTravelPlanOrderOnlyPlanning  = errors.New("plans can be reordered only while planning")
)

type TravelPlan struct {
	ID        uint `gorm:"primaryKey"`
	TravelID  uint `gorm:"not null;index"`
	StartDate *time.Time
	EndDate   *time.Time
	Place     string   `gorm:"not null"`
	Location  Location `gorm:"embedded"`
	Content   string   `gorm:"not null"`
	SortOrder int      `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
