package handler

import (
	"net/http"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TravelPlanHandler struct {
	travelPlanUsecase *usecase.TravelPlanUsecase
}

func NewTravelPlanHandler(travelPlanUsecase *usecase.TravelPlanUsecase) *TravelPlanHandler {
	return &TravelPlanHandler{
		travelPlanUsecase: travelPlanUsecase,
	}
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

	startDate, err := parseOptionalDate(req.StartDate)
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
	if err := c.ShouldBindBodyWith(&updateReq, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}
	var fields map[string]any
	if err := c.ShouldBindBodyWith(&fields, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	_, startDateSet := fields["start_date"]
	_, endDateSet := fields["end_date"]

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
		StartDate:    startDate,
		StartDateSet: startDateSet,
		EndDate:      endDate,
		EndDateSet:   endDateSet,
		Place:        updateReq.Place,
		Location:     updateReq.location(),
		Content:      updateReq.Content,
		SortOrder:    updateReq.SortOrder,
	}

	travelPlan, err := h.travelPlanUsecase.UpdateTravelPlan(userID, travelID, travelPlanID, updateInput)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTravelPlanResponse(*travelPlan))
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

	travelPlans, err := h.travelPlanUsecase.GetPublicTravelPlan(travelID)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapSlice(travelPlans, toTravelPlanResponse))
}

func (h *TravelPlanHandler) ReorderTravelPlans(c *gin.Context) {
	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid travel_id"})
		return
	}
	var req struct {
		PlanIDs *[]uint `json:"plan_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan_ids"})
		return
	}
	if err := h.travelPlanUsecase.ReorderTravelPlans(c.GetUint("userID"), travelID, *req.PlanIDs); err != nil {
		responseUsecaseError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
