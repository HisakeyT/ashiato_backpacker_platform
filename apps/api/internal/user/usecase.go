package user

import (
	"errors"

	"github.com/HisakeyT/backpacker-platform/internal/auth"
	"github.com/HisakeyT/backpacker-platform/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repository Repository
	jwtManager *auth.JWTManager
}

func NewUseCase(repository Repository, jwtManager *auth.JWTManager) *UseCase {
	return &UseCase{
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

var (
	ErrPasswordMismatch   = errors.New("password and password confirmation do not match")
	ErrEmailAlreadyUsed   = errors.New("email is already used")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (u *UseCase) Register(input RegisterInput) (*model.User, error) {
	if input.Password != input.PasswordConfirmation {
		return nil, ErrPasswordMismatch
	}

	existingUser, err := u.repository.FindByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyUsed
	}

	if !errors.Is(err, ErrUserNotFound) && err != nil {
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

func (u *UseCase) Login(input LoginInput) (*LoginOutput, error) {
	user, err := u.repository.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
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

func (u *UseCase) Me(userID uint) (*model.User, error) {
	return u.repository.FindByID(userID)
}

type UpdateMeInput struct {
	Nickname string
}

func (u *UseCase) UpdateMe(userID uint, input UpdateMeInput) (*model.User, error) {
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
