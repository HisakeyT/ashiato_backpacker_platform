package gormrepo

import (
	"errors"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/model"
)

func TestValidatePlanOrder(t *testing.T) {
	for _, tt := range []struct {
		name               string
		current, requested []uint
		wantError          bool
	}{
		{"reversed", []uint{1, 2, 3}, []uint{3, 2, 1}, false},
		{"empty travel", nil, nil, false},
		{"missing", []uint{1, 2}, []uint{1}, true},
		{"duplicate", []uint{1, 2}, []uint{1, 1}, true},
		{"foreign plan", []uint{1, 2}, []uint{1, 3}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePlanOrder(tt.current, tt.requested)
			if errors.Is(err, model.ErrInvalidTravelPlanOrder) != tt.wantError {
				t.Fatalf("error = %v", err)
			}
		})
	}
}
