package travel

import (
	"time"
)

type TravelPlanResponse struct {
	ID        uint      `json:"id"`
	TravelID  uint      `json:"travel_id"`
	StartDate string    `json:"date"`
	Place     string    `json:"place"`
	Content   string    `json:"content"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toTravelPlanResponse(p TravelPlan) TravelPlanResponse {
	return TravelPlanResponse{
		ID:        p.ID,
		TravelID:  p.TravelID,
		StartDate: p.Date.Format(dateLayout),
		Place:     p.Place,
		Content:   p.Content,
		SortOrder: p.SortOrder,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
