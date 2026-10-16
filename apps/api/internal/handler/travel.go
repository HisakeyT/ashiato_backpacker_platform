package handler

import (
	"net/http"

	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TravelHandler struct {
	travelUsecase *usecase.TravelUsecase
}

func NewTravelHandler(travelUsecase *usecase.TravelUsecase) *TravelHandler {
	return &TravelHandler{travelUsecase: travelUsecase}
}

type CreateTravelRequest struct {
	Title     string `json:"title" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	IsPublic  bool   `json:"is_public"`
}

func (h *TravelHandler) CreateTravel(c *gin.Context) {
	var req CreateTravelRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	startDate, err := parseDate(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start date format",
		})
		return
	}

	endDate, err := parseDate(req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end date format",
		})
		return
	}

	input := usecase.CreateTravelInput{
		Title:     req.Title,
		StartDate: startDate,
		EndDate:   endDate,
		IsPublic:  req.IsPublic,
	}

	travel, err := h.travelUsecase.CreateTravel(userID, input)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toTravelResponse(travel))
}

func (h *TravelHandler) GetTravels(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	travels, err := h.travelUsecase.GetTravels(userID)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapSlice(travels, toTravelResponse))
}

func (h *TravelHandler) GetTravel(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})
		return
	}

	travelDetail, err := h.travelUsecase.GetTravel(userID, travelID)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTravelDetailResponse(travelDetail))
}

type UpdateTravelRequest struct {
	Title     *string `json:"title"`
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	IsPublic  *bool   `json:"is_public"`
}

func (h *TravelHandler) UpdateTravel(c *gin.Context) {
	var updateReq UpdateTravelRequest

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})
		return
	}

	startDate, err := parseOptionalDate(updateReq.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start date format",
		})
		return
	}

	endDate, err := parseOptionalDate(updateReq.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end date format",
		})
		return
	}

	input := usecase.UpdateTravelInput{
		Title:     updateReq.Title,
		StartDate: startDate,
		EndDate:   endDate,
		IsPublic:  updateReq.IsPublic,
	}

	travel, err := h.travelUsecase.UpdateTravel(userID, travelID, input)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTravelResponse(travel))
}

func (h *TravelHandler) DeleteTravel(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})
		return
	}

	if err := h.travelUsecase.DeleteTravel(userID, travelID); err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *TravelHandler) GetPublicTravels(c *gin.Context) {
	travels, err := h.travelUsecase.GetPublicTravels()
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapSlice(travels, toPublicTravelResponse))
}

// ロジックないので、テストは書かない
func (h *TravelHandler) GetPublicTravel(c *gin.Context) {
	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})
		return
	}

	travel, err := h.travelUsecase.GetPublicTravel(travelID)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPublicTravelResponse(travel))
}

func (h *TravelHandler) CopyTravel(c *gin.Context) {
	userID := c.GetUint("userID")

	travelID, err := parseUintParam(c, "travel_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})
		return
	}

	travel, err := h.travelUsecase.CopyTravel(userID, travelID)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": travel.ID,
	})
}
