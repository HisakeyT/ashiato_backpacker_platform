package travel

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/response"
	"github.com/HisakeyT/backpacker-platform/internal/user"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase *UseCase
}

func NewHandler(useCase *UseCase) *Handler {
	return &Handler{useCase: useCase}
}

type CreateTravelRequest struct {
	Title     string `json:"title" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	IsPublic  bool   `json:"is_public"`
}

func (h *Handler) CreateTravel(c *gin.Context) {
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

	input := CreateTravelInput{
		Title:     req.Title,
		StartDate: startDate,
		EndDate:   endDate,
		IsPublic:  req.IsPublic,
	}

	travel, err := h.useCase.CreateTravel(userID, input)
	if err != nil {
		if errors.Is(err, ErrInvalidDateRange) {
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

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func (h *Handler) GetTravels(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	travels, err := h.useCase.GetTravels(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.MapSlice(travels, toTravelResponse))
}

func (h *Handler) GetTravel(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	travelIDStr := c.Param("travel_id")
	travelID64, err := strconv.ParseUint(travelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})

		return
	}

	travelID := uint(travelID64)

	travelDetail, err := h.useCase.GetTravel(userID, travelID)
	if err != nil {
		if errors.Is(err, model.ErrTravelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, user.ErrUserNotAuthorized) {
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

func (h *Handler) UpdateTravel(c *gin.Context) {
	var updateReq UpdateTravelRequest

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	travelIDStr := c.Param("travel_id")
	travelID64, err := strconv.ParseUint(travelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})

		return
	}
	travelID := uint(travelID64)

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

	input := UpdateTravelInput{
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

		if errors.Is(err, user.ErrUserNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, ErrInvalidDateRange) {
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

func parseOptionalDate(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	parsed, err := parseDate(*value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func (h *Handler) DeleteTravel(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	travelIDStr := c.Param("travel_id")
	travelID64, err := strconv.ParseUint(travelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})

		return
	}

	travelID := uint(travelID64)

	if err := h.useCase.DeleteTravel(userID, travelID); err != nil {
		if errors.Is(err, model.ErrTravelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, user.ErrUserNotAuthorized) {
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

func (h *Handler) GetPublicTravels(c *gin.Context) {
	travels, err := h.useCase.GetPublicTravels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get public travels",
		})
		return
	}

	c.JSON(http.StatusOK, response.MapSlice(travels, toPublicTravelResponse))
}

// ロジックないので、テストは書かない
func (h *Handler) GetPublicTravel(c *gin.Context) {
	travelIDStr := c.Param("travel_id")
	travelID64, err := strconv.ParseUint(travelIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid travel ID",
		})

		return
	}
	travelID := uint(travelID64)

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

func (h *Handler) CopyTravel(c *gin.Context) {
	userID := c.GetUint("userID")

	travelID64, err := strconv.ParseUint(c.Param("travel_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_id",
		})
		return
	}

	travel, err := h.useCase.CopyTravel(userID, uint(travelID64))
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
