package model

import (
	"errors"
	"time"
)

var ErrTravelNotFound = errors.New("travel not found")

type Travel struct {
	ID                 uint      `gorm:"primaryKey"`
	UserID             uint      `gorm:"not null"`
	Title              string    `gorm:"not null"`
	StartDate          time.Time `gorm:"not null"`
	EndDate            time.Time `gorm:"not null"`
	IsPublic           bool      `gorm:"not null;default:false"`
	CopiedFromTravelID *uint     `gorm:"column:copied_from_travel_id"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type CopySource struct {
	TravelID       uint
	AuthorNickname string
}

type TravelDetail struct {
	Travel     *Travel
	CopiedFrom *CopySource
}

type PublicTravel struct {
	Travel
	AuthorNickname string
}
