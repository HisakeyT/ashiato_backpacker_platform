package handler

import (
	"net/http"
	"strings"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TravelPlanHandler struct {
	travelPlanUsecase *usecase.TravelPlanUsecase
}

func NewTravelPlanHandler(travelPlanUsecase *usecase.TravelPlanUsecase) *TravelPlanHandler {
	return &TravelPlanHandler{
		travelPlanUsecase: travelPlanUsecase,
	}
}

type CreateTravelPlanRequest struct {
	StartDate string   `json:"start_date" binding:"required"`
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

func (h *TravelPlanHandler) CreateTravelPlan(c *gin.Context) {
	userID := c.GetUint("userID")

	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})
		return
	}

	var req CreateTravelPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}

	endDate, err := parseOptionalDate(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end_date",
		})
		return
	}

	input := usecase.CreateTravelPlanInput{
		StartDate: startDate,
		EndDate:   endDate,
		Place:     req.Place,
		Location: model.Location{
			Country:   req.Country,
			State:     req.State,
			City:      req.City,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		},
		Content:   req.Content,
		SortOrder: req.SortOrder,
	}

	travelPlan, err := h.travelPlanUsecase.CreateTravelPlan(
		userID,
		travelID,
		input,
	)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toTravelPlanResponse(*travelPlan))
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

func (h *TravelPlanHandler) UpdateTravelPlan(c *gin.Context) {
	userID := c.GetUint("userID")

	travelPlanID, err := parseUintParam(c, "travel_plan_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})
		return
	}

	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})
		return
	}

	var updateReq UpdateTravelPlanRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	if updateReq.Place != nil && strings.TrimSpace(*updateReq.Place) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "place must not be empty",
		})
		return
	}

	if updateReq.Content != nil && strings.TrimSpace(*updateReq.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "content must not be empty",
		})
		return
	}

	// 置き換え後
	startDate, err := parseOptionalDate(updateReq.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}

	endDate, err := parseOptionalDate(updateReq.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
		return
	}

	updateInput := usecase.UpdateTravelPlanInput{
		StartDate: startDate,
		EndDate:   endDate,
		Place:     updateReq.Place,
		Location:  updateReq.location(),
		Content:   updateReq.Content,
		SortOrder: updateReq.SortOrder,
	}

	travelPlan, err := h.travelPlanUsecase.UpdateTravelPlan(userID, travelID, travelPlanID, updateInput)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTravelPlanResponse(*travelPlan))
}

func (r UpdateTravelPlanRequest) location() *model.Location {
	if r.Country == nil && r.State == nil && r.City == nil &&
		r.Latitude == nil && r.Longitude == nil {
		return nil
	}
	return &model.Location{
		Country:   r.Country,
		State:     r.State,
		City:      r.City,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}
}

func (h *TravelPlanHandler) GetTravelPlans(c *gin.Context) {
	userID := c.GetUint("userID")

	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_id",
		})
		return
	}

	travelPlans, err := h.travelPlanUsecase.GetTravelPlans(userID, travelID)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapSlice(travelPlans, toTravelPlanResponse))
}

func (h *TravelPlanHandler) DeleteTravelPlan(c *gin.Context) {
	userID := c.GetUint("userID")

	travelPlanID, err := parseUintParam(c, "travel_plan_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_plan_id",
		})
		return
	}

	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_id",
		})
		return
	}

	err = h.travelPlanUsecase.DeleteTravelPlan(userID, travelID, travelPlanID)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TravelPlanHandler) GetPublicTravelPlan(c *gin.Context) {
	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_id",
		})
		return
	}

	travelPlans, err := h.travelPlanUsecase.GetPulicTravelPlan(uint(travelID))
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapSlice(travelPlans, toTravelPlanResponse))
}
