package usecase

import (
	"errors"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/auth"
	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUseCase_Register(t *testing.T) {
	email := "taro@example.com"

	repository := &mockrepo.MockUserRepository{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return nil, model.ErrUserNotFound
		},
		CreateFunc: func(user *model.User) error {
			user.ID = 1
			return nil
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := NewUserUseCase(repository, jwtManager)

	input := RegisterInput{
		Nickname:             "Taro",
		Email:                email,
		Password:             "password123",
		PasswordConfirmation: "password123",
	}

	user, err := useCase.Register(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID != 1 {
		t.Fatalf("expected ID 1, got %d", user.ID)
	}

	if user.Nickname != "Taro" {
		t.Fatalf("expected nickname Taro, got %s", user.Nickname)
	}

	if user.Email != email {
		t.Fatalf("unexpected email")
	}

	if user.PasswordHash == "" {
		t.Fatal("password hash should not be empty")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte("password123"),
	); err != nil {
		t.Fatalf("password was not hashed correctly: %v", err)
	}
}

func TestUseCase_Register_PasswordMismatch(t *testing.T) {
	repository := &mockrepo.MockUserRepository{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return nil, model.ErrUserNotFound
		},
		CreateFunc: func(user *model.User) error {
			t.Fatal("Create should not be called")
			return nil
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := NewUserUseCase(repository, jwtManager)

	input := RegisterInput{
		Nickname:             "Taro",
		Password:             "password123",
		PasswordConfirmation: "different-password",
	}

	_, err := useCase.Register(input)

	if !errors.Is(err, ErrPasswordMismatch) {
		t.Fatalf("expected ErrPasswordMismatch, got %v", err)
	}
}

func TestUseCase_Register_EmailAlreadyUsed(t *testing.T) {
	email := "taro@example.com"

	repository := &mockrepo.MockUserRepository{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return &model.User{
				ID:    1,
				Email: email,
			}, nil
		},
		CreateFunc: func(user *model.User) error {
			t.Fatal("Create should not be called")
			return nil
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := NewUserUseCase(repository, jwtManager)

	input := RegisterInput{
		Nickname:             "Taro",
		Email:                email,
		Password:             "password123",
		PasswordConfirmation: "password123",
	}

	_, err := useCase.Register(input)

	if !errors.Is(err, ErrEmailAlreadyUsed) {
		t.Fatalf("expected ErrEmailAlreadyUsed, got %v", err)
	}
}

func TestUseCase_Register_FindByEmailError(t *testing.T) {
	email := "taro@example.com"
	expectedErr := errors.New("database error")

	repository := &mockrepo.MockUserRepository{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return nil, expectedErr
		},
		CreateFunc: func(user *model.User) error {
			t.Fatal("Create should not be called")
			return nil
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := NewUserUseCase(repository, jwtManager)

	input := RegisterInput{
		Nickname:             "Taro",
		Email:                email,
		Password:             "password123",
		PasswordConfirmation: "password123",
	}

	_, err := useCase.Register(input)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
