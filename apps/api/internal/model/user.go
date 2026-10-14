package model

import (
	"errors"
	"time"
)

var ErrUserNotAuthorized = errors.New("user does not have permission to access this resource")
var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID           uint
	Nickname     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
