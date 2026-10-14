package user

import (
	"errors"

	"github.com/HisakeyT/backpacker-platform/internal/auth"
	"github.com/HisakeyT/backpacker-platform/internal/model"
	"testing"
)

func TestUseCase_UpdateMe(t *testing.T) {
	var updatedUser *model.User

	mockRepository := &MockRepository{
		FindByIDFunc: func(id uint) (*model.User, error) {
			return &model.User{
				ID:       1,
				Nickname: "old-nickname",
				Email:    "test@example.com",
			}, nil
		},
		UpdateFunc: func(user *model.User) error {
			updatedUser = user
			return nil
		},
	}

	jwtManager := auth.NewJWTManager("test-secret")
	useCase := NewUseCase(mockRepository, jwtManager)

	input := UpdateMeInput{
		Nickname: "new-nickname",
	}

	user, err := useCase.UpdateMe(1, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedUser == nil {
		t.Fatalf("expected updatedUser to be set, got nil")
	}

	if updatedUser.ID != 1 {
		t.Errorf("expected ID 1, got %d", user.ID)
	}

	if updatedUser.Nickname != "new-nickname" {
		t.Errorf("expected nickname new-nickname, got %s", user.Nickname)
	}

	if updatedUser.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", user.Email)
	}
}

func TestUseCase_UpdateMe_FindByIDError(t *testing.T) {
	expectedErr := errors.New("database error")

	mockRepository := &MockRepository{
		FindByIDFunc: func(id uint) (*model.User, error) {
			return nil, expectedErr
		},
	}

	useCase := NewUseCase(
		mockRepository,
		auth.NewJWTManager("test-secret"),
	)

	_, err := useCase.UpdateMe(1, UpdateMeInput{
		Nickname: "new-nickname",
	})

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestUseCase_UpdateMe_UpdateError(t *testing.T) {
	expectedErr := errors.New("database error")

	mockRepository := &MockRepository{
		FindByIDFunc: func(id uint) (*model.User, error) {
			return &model.User{
				ID:       1,
				Nickname: "old-nickname",
				Email:    "test@example.com",
			}, nil
		},
		UpdateFunc: func(user *model.User) error {
			return expectedErr
		},
	}

	useCase := NewUseCase(
		mockRepository,
		auth.NewJWTManager("test-secret"),
	)

	_, err := useCase.UpdateMe(1, UpdateMeInput{
		Nickname: "new-nickname",
	})

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}
