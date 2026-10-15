package handler

import (
	"errors"
	"net/http"

	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

// respondUsecaseError は usecase のエラーを HTTP レスポンスに変換する。
//
//	ErrTravelNotFound    -> 404
//	ErrUserNotAuthorized -> 403
//	ErrInvalidDateRange  -> 400
//	ErrInvalidCoordinates -> 400
//	ErrTravelPlanNotBelongToTravel -> 400
//	ErrTravelPlanNotFound -> 404
//	ErrEmailAlreadyUsed -> 409
//	ErrInvalidCredentials -> 401
//	ErrUserNotFound -> 404
//	ErrPasswordMismatch -> 400
//	その他                -> 500
func respondUsecaseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrTravelNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrUserNotAuthorized):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrInvalidDateRange):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrInvalidCoordinates):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrTravelPlanNotBelongToTravel):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrTravelPlanNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrEmailAlreadyUsed):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, model.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, usecase.ErrPasswordMismatch):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
