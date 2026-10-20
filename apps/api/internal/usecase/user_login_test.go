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

type stubTokenGenerator struct {
	userID uint
	err    error
}

func (generator *stubTokenGenerator) GenerateToken(userID, _ uint) (string, error) {
	generator.userID = userID
	return "test-token", generator.err
}

func TestLogin_UsesTokenGenerator(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	repository := &mock.MockUserRepository{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return &model.User{ID: 42, AuthVersion: 1, PasswordHash: string(hash)}, nil
		},
	}
	generator := &stubTokenGenerator{}
	output, err := NewUserUsecase(repository, generator).Login(LoginInput{Email: "test@example.com", Password: "password123"})
	if err != nil || output == nil || generator.userID != 42 || output.Token != "test-token" {
		t.Fatalf("output = %+v, error = %v, generated user ID = %d", output, err, generator.userID)
	}

	generator.err = errors.New("token error")
	output, err = NewUserUsecase(repository, generator).Login(LoginInput{Email: "test@example.com", Password: "password123"})
	if !errors.Is(err, generator.err) || output != nil {
		t.Fatalf("output = %+v, error = %v", output, err)
	}
}

func TestLogin_Success(t *testing.T) {
	email := "test@example.com"

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	mockRepo := &mock.MockUserRepository{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return &model.User{
				ID:           1,
				AuthVersion:  1,
				Nickname:     "Taro",
				Email:        "test@example.com",
				PasswordHash: string(hashedPassword),
			}, nil
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := NewUserUsecase(mockRepo, jwtManager)

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
		func(token *jwt.Token) (any, error) {
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
	mockRepo := &mock.MockUserRepository{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return nil, model.ErrUserNotFound
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := NewUserUsecase(mockRepo, jwtManager)

	_, err := useCase.Login(LoginInput{
		Email:    "notfound@example.com",
		Password: "password123",
	})

	if !errors.Is(err, model.ErrInvalidCredentials) {
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

	mockRepo := &mock.MockUserRepository{
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
	useCase := NewUserUsecase(mockRepo, jwtManager)

	_, err = useCase.Login(LoginInput{
		Email:    "test@example.com",
		Password: "wrong-password",
	})

	if !errors.Is(err, model.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
