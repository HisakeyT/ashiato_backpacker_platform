package travel

import (
	"time"
)

type TravelPlanResponse struct {
	ID        uint      `json:"id"`
	TravelID  uint      `json:"travel_id"`
	StartDate string    `json:"start_date"`
	EndDate   *string   `json:"end_date"`
	Place     string    `json:"place"`
	Content   string    `json:"content"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toTravelPlanResponse(p TravelPlan) TravelPlanResponse {
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
		Content:   p.Content,
		SortOrder: p.SortOrder,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
