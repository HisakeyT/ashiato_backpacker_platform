package usecase

import (
	"errors"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/auth"
	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin_Success(t *testing.T) {
	email := "test@example.com"

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	mockRepo := &mockrepo.MockUserRepository{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return &model.User{
				ID:           1,
				Nickname:     "Taro",
				Email:        "test@example.com",
				PasswordHash: string(hashedPassword),
			}, nil
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := NewUserUseCase(mockRepo, jwtManager)

	result, err := useCase.Login(LoginInput{
		Email:    email,
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.User.ID != 1 {
		t.Errorf("expected user ID 1, got %d", result.User.ID)
	}

	if result.User.Email != email {
		t.Errorf("expected email %s, got %v", email, result.User.Email)
	}

	if result.Token == "" {
		t.Errorf("expected token to be generated")
	}

	token, err := jwt.Parse(
		result.Token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte("test-secret"), nil
		},
	)

	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if !token.Valid {
		t.Error("expected token to be valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected claims to be jwt.MapClaims")
	}

	if claims["user_id"] != float64(1) {
		t.Errorf("expected user_id 1, got %v", claims["user_id"])
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := &mockrepo.MockUserRepository{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return nil, model.ErrUserNotFound
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := NewUserUseCase(mockRepo, jwtManager)

	_, err := useCase.Login(LoginInput{
		Email:    "notfound@example.com",
		Password: "password123",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	mockRepo := &mockrepo.MockUserRepository{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return &model.User{
				ID:           1,
				Nickname:     "Taro",
				Email:        "test@example.com",
				PasswordHash: string(hashedPassword),
			}, nil
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := NewUserUseCase(mockRepo, jwtManager)

	_, err = useCase.Login(LoginInput{
		Email:    "test@example.com",
		Password: "wrong-password",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
