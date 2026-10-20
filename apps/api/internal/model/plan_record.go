package model

import (
	"errors"
	"time"
)

var (
	ErrPlanRecordNotFound = errors.New("plan record not found")
	ErrInvalidPlanRecord  = errors.New("invalid plan record")
)

type RecordKind string

const (
	RecordSightseeing RecordKind = "sightseeings"
	RecordStay        RecordKind = "stays"
	RecordMeal        RecordKind = "meals"
	RecordOther       RecordKind = "others"
)

func (k RecordKind) Valid() bool {
	switch k {
	case RecordSightseeing, RecordStay, RecordMeal, RecordOther:
		return true
	default:
		return false
	}
}

type PlanRecord struct {
	ID           uint `gorm:"primaryKey"`
	TravelPlanID uint
	Name         string
	MealType     *string
	LocationText *string
	MapURL       *string
	Cost         *uint
	SortOrder    int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
