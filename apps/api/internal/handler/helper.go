package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

const dateLayout = "2006-01-02"

func parseDate(s string) (time.Time, error) {
	return time.Parse(dateLayout, s)
}

func parseOptionalDate(v *string) (*time.Time, error) {
	if v == nil || *v == "" {
		return nil, nil
	}
	t, err := parseDate(*v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseUintParam(c *gin.Context, name string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func mapSlice[T any, R any](items []T, convert func(T) R) []R {
	result := make([]R, 0, len(items))

	for _, item := range items {
		result = append(result, convert(item))
	}

	return result
}
