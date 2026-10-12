package travel

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/HisakeyT/backpacker-platform/internal/response"
	"github.com/HisakeyT/backpacker-platform/internal/user"
	"github.com/gin-gonic/gin"
)

type PlanHandler struct {
	useCase *PlanUseCase
}

func NewPlanHandler(useCase *PlanUseCase) *PlanHandler {
	return &PlanHandler{
		useCase: useCase,
	}
}

type CreateTravelPlanRequest struct {
	StartDate string `json:"date" binding:"required"`
	Place     string `json:"place" binding:"required"`
	Content   string `json:"content" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

func (h *PlanHandler) CreateTravelPlan(c *gin.Context) {
	userID := c.GetUint("userID")

	travelID, err := strconv.ParseUint(c.Param("travel_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_id",
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

	date, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date",
		})
		return
	}

	input := CreateTravelPlanInput{
		StartDate: date,
		Place:     req.Place,
		Content:   req.Content,
		SortOrder: req.SortOrder,
	}

	travelPlan, err := h.useCase.CreateTravelPlan(
		userID,
		uint(travelID),
		input,
	)

	if err != nil {
		if errors.Is(err, ErrTravelNotFound) {
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

	c.JSON(http.StatusCreated, toTravelPlanResponse(*travelPlan))
}

type UpdateTravelPlanRequest struct {
	StartDate *string `json:"date"`
	Place     *string `json:"place"`
	Content   *string `json:"content"`
	SortOrder *int    `json:"sort_order"`
}

func (h *PlanHandler) UpdateTravelPlan(c *gin.Context) {
	userID := c.GetUint("userID")

	travelPlanID64, err := strconv.ParseUint(c.Param("travel_plan_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_plan_id",
		})
		return
	}

	travelID64, err := strconv.ParseUint(c.Param("travel_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_id",
		})
		return
	}

	travelID := uint(travelID64)
	travelPlanID := uint(travelPlanID64)

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

	var date *time.Time
	if updateReq.StartDate != nil {
		parsed, err := time.Parse("2006-01-02", *updateReq.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid date",
			})
			return
		}

		date = &parsed
	}

	updateInput := UpdateTravelPlanInput{
		StartDate: date,
		Place:     updateReq.Place,
		Content:   updateReq.Content,
		SortOrder: updateReq.SortOrder,
	}

	travelPlan, err := h.useCase.UpdateTravelPlan(userID, travelID, travelPlanID, updateInput)
	if err != nil {
		if errors.Is(err, ErrTravelNotFound) {
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

		if errors.Is(err, ErrTravelPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, ErrTravelPlanNotBelongToTravel) {
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

	c.JSON(http.StatusOK, toTravelPlanResponse(*travelPlan))
}

func (h *PlanHandler) GetTravelPlans(c *gin.Context) {
	userID := c.GetUint("userID")

	travelID, err := strconv.ParseUint(c.Param("travel_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_id",
		})
		return
	}

	travelPlans, err := h.useCase.GetTravelPlans(userID, uint(travelID))
	if err != nil {
		if errors.Is(err, ErrTravelNotFound) {
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

	c.JSON(http.StatusOK, response.MapSlice(travelPlans, toTravelPlanResponse))
}

func (h *PlanHandler) DeleteTravelPlan(c *gin.Context) {
	userID := c.GetUint("userID")

	travelPlanID64, err := strconv.ParseUint(c.Param("travel_plan_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_plan_id",
		})
		return
	}

	travelID64, err := strconv.ParseUint(c.Param("travel_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_id",
		})
		return
	}

	err = h.useCase.DeleteTravelPlan(userID, uint(travelID64), uint(travelPlanID64))
	if err != nil {
		if errors.Is(err, ErrTravelNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, user.ErrUserNotAuthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, ErrTravelPlanNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, ErrTravelPlanNotBelongToTravel) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PlanHandler) GetPublicTravelPlan(c *gin.Context) {
	travelID, err := strconv.ParseUint(c.Param("travel_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid travel_id",
		})
		return
	}

	travelPlans, err := h.useCase.GetPulicTravelPlan(uint(travelID))
	if err != nil {
		if errors.Is(err, ErrTravelNotFound) {
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

	c.JSON(http.StatusOK, response.MapSlice(travelPlans, toTravelPlanResponse))
}
