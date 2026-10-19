package handler

import (
	"errors"
	"net/http"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/gin-gonic/gin"
)

type errorStatus struct {
	err    error
	status int
}

var errorStatuses = []errorStatus{
	// 400
	{model.ErrInvalidDateRange, http.StatusBadRequest},
	{model.ErrTravelPlanOutsideTravelDates, http.StatusBadRequest},
	{model.ErrPasswordMismatch, http.StatusBadRequest},
	{model.ErrInvalidCoordinates, http.StatusBadRequest},
	{model.ErrInvalidTravelPlanText, http.StatusBadRequest},
	{model.ErrTravelPlanNotBelongToTravel, http.StatusBadRequest},
	{model.ErrInvalidScheduleMode, http.StatusBadRequest},
	{model.ErrPlanStartDateRequired, http.StatusBadRequest},
	{model.ErrTravelPlanStartDateRequired, http.StatusBadRequest},
	{model.ErrTravelPlanEndWithoutStart, http.StatusBadRequest},
	{model.ErrInvalidTravelPlanOrder, http.StatusBadRequest},
	{model.ErrTravelPlanOrderOnlyPlanning, http.StatusBadRequest},

	// 401
	{model.ErrInvalidCredentials, http.StatusUnauthorized},

	// 403
	{model.ErrUserNotAuthorized, http.StatusForbidden},

	// 404
	{model.ErrTravelNotFound, http.StatusNotFound},
	{model.ErrTravelPlanNotFound, http.StatusNotFound},
	{model.ErrUserNotFound, http.StatusNotFound},

	// 409
	{model.ErrEmailAlreadyUsed, http.StatusConflict},
}

// responseUsecaseError は usecase のエラーを HTTP レスポンスに変換する。
// 対応は errorStatuses を参照。未登録のエラーは 500。
func responseUsecaseError(c *gin.Context, err error) {
	for _, e := range errorStatuses {
		if errors.Is(err, e.err) {
			c.JSON(e.status, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
