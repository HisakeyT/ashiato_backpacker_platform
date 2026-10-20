package auth

import (
	"math"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret")
	userID := uint(1)

	tokenString, err := manager.GenerateToken(userID, 1)

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

func TestJWTExpiresAt24HoursAndRequiresClaims(t *testing.T) {
	fixed := time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC)
	manager := NewJWTManager("test-secret")
	manager.now = func() time.Time { return fixed }
	token, err := manager.GenerateToken(7, 3)
	require.NoError(t, err)
	identity, err := manager.ValidateIdentity(token)
	require.NoError(t, err)
	assert.Equal(t, TokenIdentity{UserID: 7, AuthVersion: 3}, identity)

	manager.now = func() time.Time { return fixed.Add(24*time.Hour - time.Second) }
	_, err = manager.ValidateIdentity(token)
	require.NoError(t, err)
	manager.now = func() time.Time { return fixed.Add(24 * time.Hour) }
	_, err = manager.ValidateIdentity(token)
	require.Error(t, err)

	manager.now = func() time.Time { return fixed }
	for _, claims := range []jwt.MapClaims{
		{"user_id": 7, "auth_version": 3, "iat": fixed.Unix(), "exprires_at": fixed.Add(24 * time.Hour).Unix()},
		{"user_id": 7, "auth_version": 3, "iat": fixed.Unix()},
		{"user_id": 7, "exp": fixed.Add(time.Hour).Unix(), "iat": fixed.Unix()},
		{"user_id": 7, "auth_version": 3, "exp": fixed.Add(time.Hour).Unix()},
		{"user_id": -1, "auth_version": 3, "exp": fixed.Add(time.Hour).Unix(), "iat": fixed.Unix()},
		{"user_id": 1.5, "auth_version": 3, "exp": fixed.Add(time.Hour).Unix(), "iat": fixed.Unix()},
		{"user_id": math.MaxFloat64, "auth_version": 3, "exp": fixed.Add(time.Hour).Unix(), "iat": fixed.Unix()},
	} {
		old, signErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))
		require.NoError(t, signErr)
		_, err = manager.ValidateIdentity(old)
		require.Error(t, err, "claims: %v", claims)
	}
}

func TestJWTRejectsOtherSigningMethodsAndTampering(t *testing.T) {
	manager := NewJWTManager("test-secret")
	valid, err := manager.GenerateToken(7, 1)
	require.NoError(t, err)
	_, err = manager.ValidateIdentity(valid + "x")
	require.Error(t, err)
	claims := jwt.MapClaims{"user_id": 7, "auth_version": 1, "iat": time.Now().Unix(), "exp": time.Now().Add(time.Hour).Unix()}
	otherMethod, err := jwt.NewWithClaims(jwt.SigningMethodHS384, claims).SignedString([]byte("test-secret"))
	require.NoError(t, err)
	_, err = manager.ValidateIdentity(otherMethod)
	require.Error(t, err)
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret")

	token, err := manager.GenerateToken(123, 1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	identity, err := manager.ValidateIdentity(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if identity.UserID != 123 || identity.AuthVersion != 1 {
		t.Errorf("unexpected identity: %+v", identity)
	}
}

func TestJWTManager_ValidateToken_InvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret")

	_, err := manager.ValidateIdentity("invalid-token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJWTManager_ValidateToken_WrongSecret(t *testing.T) {
	manager := NewJWTManager("test-secret")
	otherManager := NewJWTManager("other-secret")

	token, err := otherManager.GenerateToken(123, 1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = manager.ValidateIdentity(token)
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

	_, err = manager.ValidateIdentity(tokenString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
