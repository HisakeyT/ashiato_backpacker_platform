package user

import (
	"errors"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/auth"
	"github.com/HisakeyT/backpacker-platform/internal/model"
)

func TestUseCase_Me(t *testing.T) {
	repository := &fakeRepository{
		user: &model.User{
			ID:       1,
			Nickname: "test",
			Email:    "test@example.com",
		},
	}

	jwtManager := auth.NewJWTManager("secret")
	useCase := NewUseCase(repository, jwtManager)

	user, err := useCase.Me(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user ID %d, got %d", 1, user.ID)
	}
}

func TestUseCase_Me_UserNotFound(t *testing.T) {
	repository := &fakeRepository{err: ErrUserNotFound}
	jwtManager := auth.NewJWTManager("secret")
	useCase := NewUseCase(repository, jwtManager)

	user, err := useCase.Me(999)

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}
}
