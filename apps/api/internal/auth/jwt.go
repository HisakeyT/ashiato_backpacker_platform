package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey string
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
	}
}

func (j *JWTManager) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     userID,
		"exprires_at": time.Now().Add(time.Hour * 24).Unix(),
		"issued_at":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(j.secretKey))
}

func (m *JWTManager) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(m.secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user_id")
	}

	return uint(userIDFloat), nil
}
