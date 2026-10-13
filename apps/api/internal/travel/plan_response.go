package travel

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
	"time"
)

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
