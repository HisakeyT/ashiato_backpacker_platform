package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret")
	userID := uint(1)

	tokenString, err := manager.GenerateToken(userID)

	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			return []byte("test-secret"), nil
		},
	)

	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, float64(1), claims["user_id"])
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret")

	token, err := manager.GenerateToken(123)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	userID, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if userID != 123 {
		t.Errorf("expected userID 123, got %d", userID)
	}
}

func TestJWTManager_ValidateToken_InvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret")

	_, err := manager.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJWTManager_ValidateToken_WrongSecret(t *testing.T) {
	manager := NewJWTManager("test-secret")
	otherManager := NewJWTManager("other-secret")

	token, err := otherManager.GenerateToken(123)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = manager.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJWTManager_ValidateToken_ExpiredToken(t *testing.T) {
	manager := NewJWTManager("test-secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uint(123),
		"exp":     time.Now().Add(-time.Hour).Unix(),
		"iat":     time.Now().Add(-2 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = manager.ValidateToken(tokenString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
