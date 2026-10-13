package model

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidCoordinates = errors.New("invalid coordinates")

type Location struct {
	Country   *string
	State     *string
	City      *string
	Latitude  *float64
	Longitude *float64
}

func (l *Location) Normalize() {
	l.Country = normalizeOptional(l.Country)
	l.State = normalizeOptional(l.State)
	l.City = normalizeOptional(l.City)
}

func (l Location) Validate() error {
	if (l.Latitude == nil) != (l.Longitude == nil) {
		return fmt.Errorf("%w: latitude and longitude must be set together", ErrInvalidCoordinates)
	}
	if l.Latitude != nil && (*l.Latitude < -90 || *l.Latitude > 90) {
		return fmt.Errorf("%w: latitude must be between -90 and 90", ErrInvalidCoordinates)
	}
	if l.Longitude != nil && (*l.Longitude < -180 || *l.Longitude > 180) {
		return fmt.Errorf("%w: longitude must be between -180 and 180", ErrInvalidCoordinates)
	}
	return nil
}

func normalizeOptional(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	return &t
}
