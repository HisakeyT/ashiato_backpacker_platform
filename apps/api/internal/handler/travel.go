package handler

import (
	"errors"
	"net/http"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TravelHandler struct {
	useCase *usecase.TravelUsecase
}

func NewTravelHandler(useCase *usecase.TravelUsecase) *TravelHandler {
	return &TravelHandler{useCase: useCase}
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

	travel, err := h.useCase.CreateTravel(userID, input)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, toTravelResponse(travel))
}

func (h *TravelHandler) GetTravels(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	travels, err := h.useCase.GetTravels(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
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

	travelDetail, err := h.useCase.GetTravel(userID, travelID)
	if err != nil {
		if errors.Is(err, model.ErrTravelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, model.ErrUserNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

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

	travel, err := h.useCase.UpdateTravel(userID, travelID, input)
	if err != nil {
		if errors.Is(err, model.ErrTravelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, model.ErrUserNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, usecase.ErrInvalidDateRange) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

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

	if err := h.useCase.DeleteTravel(userID, travelID); err != nil {
		if errors.Is(err, model.ErrTravelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, model.ErrUserNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *TravelHandler) GetPublicTravels(c *gin.Context) {
	travels, err := h.useCase.GetPublicTravels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get public travels",
		})
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

	travel, err := h.useCase.GetPublicTravel(travelID)
	if err != nil {
		if errors.Is(err, model.ErrTravelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
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

	travel, err := h.useCase.CopyTravel(userID, travelID)
	if err != nil {
		if errors.Is(err, model.ErrTravelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": travel.ID,
	})
}
