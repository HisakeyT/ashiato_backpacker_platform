package model

import (
	"errors"
	"time"
)

var ErrTravelNotFound = errors.New("travel not found")
var ErrInvalidDateRange = errors.New("start date must not be after end date")

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

	// 複製時など、明示的に Preload したときだけ入る。通常は nil。
	Plans []TravelPlan `gorm:"foreignKey:TravelID;<-:false" json:"-"`
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

func (s Travel) Copy(newUserID uint) *Travel {
	originalID := s.ID

	plans := make([]TravelPlan, len(s.Plans))
	for i, p := range s.Plans {
		var endDate *time.Time
		if p.EndDate != nil {
			e := *p.EndDate
			endDate = &e
		}
		plans[i] = TravelPlan{
			StartDate: p.StartDate,
			EndDate:   endDate,
			Place:     p.Place,
			Location:  p.Location,
			Content:   p.Content,
			SortOrder: p.SortOrder,
		}
	}

	return &Travel{
		UserID:             newUserID,
		Title:              s.Title,
		StartDate:          s.StartDate,
		EndDate:            s.EndDate,
		IsPublic:           false, // 新しい旅行は非公開で作成される
		CopiedFromTravelID: &originalID,
		Plans:              plans,
	}
}
