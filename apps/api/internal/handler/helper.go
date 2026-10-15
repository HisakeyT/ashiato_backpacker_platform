package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

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
