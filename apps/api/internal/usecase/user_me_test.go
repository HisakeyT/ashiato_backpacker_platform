package usecase

import (
	"errors"
	"testing"

	"github.com/HisakeyT/backpacker-platform/internal/auth"
	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository/mock"
)

func TestUserUsecase_Me(t *testing.T) {
	repo := &mockrepo.MockUserRepository{
		FindByIDFunc: func(id uint) (*model.User, error) {
			return &model.User{
				ID:       1,
				Nickname: "test",
				Email:    "test@example.com",
			}, nil
		},
	}

	jwtManager := auth.NewJWTManager("secret")
	uc := NewUserUseCase(repo, jwtManager)

	user, err := uc.Me(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != 1 {
		t.Errorf("expected user ID %d, got %d", 1, user.ID)
	}
}

func TestUserUsecase_Me_UserNotFound(t *testing.T) {
	repo := &mockrepo.MockUserRepository{
		FindByIDFunc: func(id uint) (*model.User, error) {
			return nil, model.ErrUserNotFound
		},
	}

	jwtManager := auth.NewJWTManager("secret")
	uc := NewUserUseCase(repo, jwtManager)

	user, err := uc.Me(999)

	if !errors.Is(err, model.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	if user != nil {
		t.Errorf("expected nil user, got %+v", user)
	}
}
