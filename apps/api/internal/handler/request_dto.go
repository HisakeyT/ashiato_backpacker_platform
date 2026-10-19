package handler

import (
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

// user
type RegisterRequest struct {
	Nickname             string `json:"nickname" binding:"required"`
	Email                string `json:"email" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

// travel
type CreateTravelRequest struct {
	Title        string              `json:"title" binding:"required"`
	StartDate    string              `json:"start_date" binding:"required"`
	EndDate      string              `json:"end_date" binding:"required"`
	IsPublic     bool                `json:"is_public"`
	ScheduleMode *model.ScheduleMode `json:"schedule_mode"`
}

type UpdateTravelRequest struct {
	Title        *string             `json:"title"`
	StartDate    *string             `json:"start_date"`
	EndDate      *string             `json:"end_date"`
	IsPublic     *bool               `json:"is_public"`
	ScheduleMode *model.ScheduleMode `json:"schedule_mode"`
}

// travel_plan
type CreateTravelPlanRequest struct {
	StartDate *string  `json:"start_date"`
	EndDate   *string  `json:"end_date"`
	Place     string   `json:"place" binding:"required"`
	Country   *string  `json:"country"`
	State     *string  `json:"state"`
	City      *string  `json:"city"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Content   string   `json:"content" binding:"required"`
	SortOrder int      `json:"sort_order"`
}

type UpdateTravelPlanRequest struct {
	StartDate *string  `json:"start_date"`
	EndDate   *string  `json:"end_date"`
	Place     *string  `json:"place"`
	Country   *string  `json:"country"`
	State     *string  `json:"state"`
	City      *string  `json:"city"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Content   *string  `json:"content"`
	SortOrder *int     `json:"sort_order"`
}

func (r UpdateTravelPlanRequest) location() *model.Location {
	return &model.Location{
		Country:   r.Country,
		State:     r.State,
		City:      r.City,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}
}
