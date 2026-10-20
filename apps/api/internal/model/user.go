package model

import (
	"errors"
	"time"
)

var (
	ErrUserNotAuthorized  = errors.New("user does not have permission to access this resource")
	ErrUserNotFound       = errors.New("user not found")
	ErrPasswordMismatch   = errors.New("password and password confirmation do not match")
	ErrEmailAlreadyUsed   = errors.New("email is already used")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type User struct {
	ID           uint
	Nickname     string
	Email        string
	PasswordHash string
	AuthVersion  uint `gorm:"not null;default:1"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
