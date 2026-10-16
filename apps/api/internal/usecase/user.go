package usecase

import (
	"errors"

	"github.com/HisakeyT/backpacker-platform/internal/auth"
	"github.com/HisakeyT/backpacker-platform/internal/model"
	"github.com/HisakeyT/backpacker-platform/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repository repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewUserUsecase(repository repository.UserRepository, jwtManager *auth.JWTManager) *UserUsecase {
	return &UserUsecase{
		repository: repository,
		jwtManager: jwtManager,
	}
}

type RegisterInput struct {
	Nickname             string
	Email                string
	Password             string
	PasswordConfirmation string
}

func (u *UserUsecase) Register(input RegisterInput) (*model.User, error) {
	if input.Password != input.PasswordConfirmation {
		return nil, model.ErrPasswordMismatch
	}

	existingUser, err := u.repository.FindByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, model.ErrEmailAlreadyUsed
	}

	if !errors.Is(err, model.ErrUserNotFound) && err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Nickname:     input.Nickname,
		Email:        input.Email,
		PasswordHash: string(passwordHash),
	}

	if err := u.repository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	User  *model.User
	Token string
}

func (u *UserUsecase) Login(input LoginInput) (*LoginOutput, error) {
	user, err := u.repository.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			return nil, model.ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, model.ErrInvalidCredentials
	}

	token, err := u.jwtManager.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		User:  user,
		Token: token,
	}, nil
}

func (u *UserUsecase) Me(userID uint) (*model.User, error) {
	return u.repository.FindByID(userID)
}

type UpdateMeInput struct {
	Nickname string
}

func (u *UserUsecase) UpdateMe(userID uint, input UpdateMeInput) (*model.User, error) {
	user, err := u.repository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	user.Nickname = input.Nickname

	if err := u.repository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
