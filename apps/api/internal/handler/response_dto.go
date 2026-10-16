package handler

import (
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
)

// user
type UserResponse struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// travel
// TravelResponse は自分の旅行の基本情報。
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

// PublicTravelResponse は公開旅行の一覧・詳細用。user_id は出さず、投稿者名を返す。
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

// CopiedFromResponse はコピー元の旅行 ID と投稿者名。
type CopiedFromResponse struct {
	TravelID       uint   `json:"travel_id"`
	AuthorNickname string `json:"author_nickname"`
}

// TravelDetailResponse は旅行詳細。コピー元がある場合のみ CopiedFrom を持つ。
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

// travel plan
// TravelPlanResponse は旅行内の1件の予定。EndDate と位置情報は任意。
type TravelPlanResponse struct {
	ID        uint      `json:"id"`
	TravelID  uint      `json:"travel_id"`
	StartDate string    `json:"start_date"`
	EndDate   *string   `json:"end_date"`
	Place     string    `json:"place"`
	Country   *string   `json:"country"`
	State     *string   `json:"state"`
	City      *string   `json:"city"`
	Latitude  *float64  `json:"latitude"`
	Longitude *float64  `json:"longitude"`
	Content   string    `json:"content"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toTravelPlanResponse(p model.TravelPlan) TravelPlanResponse {
	var endDate *string
	if p.EndDate != nil {
		formattedEndDate := p.EndDate.Format(dateLayout)
		endDate = &formattedEndDate
	}

	return TravelPlanResponse{
		ID:        p.ID,
		TravelID:  p.TravelID,
		StartDate: p.StartDate.Format(dateLayout),
		EndDate:   endDate,
		Place:     p.Place,
		Country:   p.Location.Country,
		State:     p.Location.State,
		City:      p.Location.City,
		Latitude:  p.Location.Latitude,
		Longitude: p.Location.Longitude,
		Content:   p.Content,
		SortOrder: p.SortOrder,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
