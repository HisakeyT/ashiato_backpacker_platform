package handler

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
	"time"
)

const dateLayout = "2006-01-02"

type TravelResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CopiedFromResponse struct {
	TravelID       uint   `json:"travel_id"`
	AuthorNickname string `json:"author_nickname"`
}

type TravelDetailResponse struct {
	TravelResponse
	CopiedFrom *CopiedFromResponse `json:"copied_from"`
}

func toTravelDetailResponse(detail *model.TravelDetail) TravelDetailResponse {
	res := TravelDetailResponse{
		TravelResponse: toTravelResponse(detail.Travel),
	}

	if detail.CopiedFrom != nil {
		res.CopiedFrom = &CopiedFromResponse{
			TravelID:       detail.CopiedFrom.TravelID,
			AuthorNickname: detail.CopiedFrom.AuthorNickname,
		}
	}

	return res
}

func toTravelResponse(t *model.Travel) TravelResponse {
	return TravelResponse{
		ID:        t.ID,
		UserID:    t.UserID,
		Title:     t.Title,
		StartDate: t.StartDate.Format(dateLayout),
		EndDate:   t.EndDate.Format(dateLayout),
		IsPublic:  t.IsPublic,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

type PublicTravelResponse struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date"`
	AuthorNickname string    `json:"author_nickname"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func toPublicTravelResponse(t *model.PublicTravel) PublicTravelResponse {
	return PublicTravelResponse{
		ID:             t.ID,
		Title:          t.Title,
		StartDate:      t.StartDate.Format(dateLayout),
		EndDate:        t.EndDate.Format(dateLayout),
		AuthorNickname: t.AuthorNickname,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}
