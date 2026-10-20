package handler

import (
	"encoding/json"
	"net/http"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

type PlanRecordHandler struct{ usecase *usecase.PlanRecordUsecase }

func NewPlanRecordHandler(uc *usecase.PlanRecordUsecase) *PlanRecordHandler {
	return &PlanRecordHandler{usecase: uc}
}

type recordParams struct {
	travelID, planID, recordID uint
	kind                       model.RecordKind
}

func parseRecordParams(c *gin.Context, withRecord bool) (recordParams, bool) {
	var p recordParams
	var err error
	p.travelID, err = parseUintParam(c, "travel_id")
	if err != nil || p.travelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid travel_id"})
		return p, false
	}
	p.planID, err = parseUintParam(c, "travel_plan_id")
	if err != nil || p.planID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid travel_plan_id"})
		return p, false
	}
	p.kind = model.RecordKind(c.Param("kind"))
	if !p.kind.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record kind"})
		return p, false
	}
	if withRecord {
		p.recordID, err = parseUintParam(c, "record_id")
		if err != nil || p.recordID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid record_id"})
			return p, false
		}
	}
	return p, true
}

func recordResponse(record model.PlanRecord) gin.H {
	return gin.H{
		"id": record.ID, "travel_plan_id": record.TravelPlanID, "name": record.Name,
		"meal_type": record.MealType, "location_text": record.LocationText, "map_url": record.MapURL,
		"cost": record.Cost, "sort_order": record.SortOrder,
	}
}

func readRecordBody(c *gin.Context, kind model.RecordKind, create bool) (usecase.RecordInput, usecase.RecordPatch, bool) {
	var fields map[string]json.RawMessage
	if err := c.ShouldBindJSON(&fields); err != nil || fields == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return usecase.RecordInput{}, usecase.RecordPatch{}, false
	}
	input := usecase.RecordInput{}
	patch := usecase.RecordPatch{}
	for key, raw := range fields {
		var value *string
		if err := json.Unmarshal(raw, &value); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return input, patch, false
		}
		switch key {
		case "name":
			if value == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
				return input, patch, false
			}
			input.Name, patch.Name = *value, value
		case "location_text":
			input.LocationText, patch.LocationText = value, &value
		case "map_url":
			input.MapURL, patch.MapURL = value, &value
		case "meal_type":
			if kind != model.RecordMeal {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid meal_type"})
				return input, patch, false
			}
			input.MealType, patch.MealType = value, &value
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown field"})
			return input, patch, false
		}
	}
	if create && patch.Name == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return input, patch, false
	}
	return input, patch, true
}

func (h *PlanRecordHandler) List(c *gin.Context) {
	h.list(c, false)
}

func (h *PlanRecordHandler) ListPublic(c *gin.Context) {
	h.list(c, true)
}

func (h *PlanRecordHandler) list(c *gin.Context, public bool) {
	p, ok := parseRecordParams(c, false)
	if !ok {
		return
	}
	records, err := h.usecase.List(c.GetUint("userID"), p.travelID, p.planID, p.kind, public)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, mapSlice(records, recordResponse))
}

func (h *PlanRecordHandler) Create(c *gin.Context) {
	p, ok := parseRecordParams(c, false)
	if !ok {
		return
	}
	input, _, ok := readRecordBody(c, p.kind, true)
	if !ok {
		return
	}
	record, err := h.usecase.Create(c.GetUint("userID"), p.travelID, p.planID, p.kind, input)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusCreated, recordResponse(*record))
}

func (h *PlanRecordHandler) Update(c *gin.Context) {
	p, ok := parseRecordParams(c, true)
	if !ok {
		return
	}
	_, patch, ok := readRecordBody(c, p.kind, false)
	if !ok {
		return
	}
	record, err := h.usecase.Update(c.GetUint("userID"), p.travelID, p.planID, p.recordID, p.kind, patch)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}
	c.JSON(http.StatusOK, recordResponse(*record))
}

func (h *PlanRecordHandler) Delete(c *gin.Context) {
	p, ok := parseRecordParams(c, true)
	if !ok {
		return
	}
	if err := h.usecase.Delete(c.GetUint("userID"), p.travelID, p.planID, p.recordID, p.kind); err != nil {
		responseUsecaseError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
