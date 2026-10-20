package auth

import (
	"fmt"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenLifetime = 24 * time.Hour
const maxExactJSONInteger = 1<<53 - 1

type JWTManager struct {
	secretKey string
	now       func() time.Time
}

type TokenIdentity struct {
	UserID      uint
	AuthVersion uint
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{secretKey: secretKey, now: time.Now}
}

func (j *JWTManager) GenerateToken(userID, authVersion uint) (string, error) {
	if userID == 0 || authVersion == 0 || uint64(userID) > maxExactJSONInteger ||
		uint64(authVersion) > maxExactJSONInteger {
		return "", fmt.Errorf("invalid token identity")
	}
	now := j.now()
	claims := jwt.MapClaims{
		"user_id":      userID,
		"auth_version": authVersion,
		"exp":          now.Add(tokenLifetime).Unix(),
		"iat":          now.Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.secretKey))
}

func positiveUintClaim(claims jwt.MapClaims, key string) (uint, error) {
	value, ok := claims[key].(float64)
	if !ok || math.IsNaN(value) || math.IsInf(value, 0) || value < 1 ||
		value > maxExactJSONInteger || math.Trunc(value) != value {
		return 0, fmt.Errorf("invalid %s", key)
	}
	return uint(value), nil
}

func (j *JWTManager) ValidateIdentity(tokenString string) (TokenIdentity, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(j.now),
	)
	token, err := parser.Parse(tokenString, func(*jwt.Token) (any, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return TokenIdentity{}, fmt.Errorf("invalid token: %w", err)
	}
	if !token.Valid {
		return TokenIdentity{}, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return TokenIdentity{}, fmt.Errorf("invalid claims")
	}
	if _, err := positiveUintClaim(claims, "iat"); err != nil {
		return TokenIdentity{}, err
	}
	userID, err := positiveUintClaim(claims, "user_id")
	if err != nil {
		return TokenIdentity{}, err
	}
	version, err := positiveUintClaim(claims, "auth_version")
	if err != nil {
		return TokenIdentity{}, err
	}
	return TokenIdentity{UserID: userID, AuthVersion: version}, nil
}
