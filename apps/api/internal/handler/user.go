package handler

import (
	"net/http"

	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase *usecase.UserUsecase
}

func NewHandler(useCase *usecase.UserUsecase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type RegisterRequest struct {
	Nickname             string `json:"nickname" binding:"required"`
	Email                string `json:"email" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	input := usecase.RegisterInput{
		Nickname:             req.Nickname,
		Email:                req.Email,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}

	u, err := h.useCase.Register(input)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, UserResponse{
		ID:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var input usecase.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	result, err := h.useCase.Login(input)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: result.Token,
		User: UserResponse{
			ID:       result.User.ID,
			Nickname: result.User.Nickname,
			Email:    result.User.Email,
		},
	})
}

func (h *Handler) Me(c *gin.Context) {
	userID, _ := c.Get("userID")

	user, err := h.useCase.Me(userID.(uint))
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Email:    user.Email,
	})
}

func (h *Handler) UpdateMe(c *gin.Context) {
	var input usecase.UpdateMeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	userID, _ := c.Get("userID")

	user, err := h.useCase.UpdateMe(userID.(uint), input)
	if err != nil {
		responseUsecaseError(c, err)
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:       user.ID,
		Nickname: user.Nickname,
		Email:    user.Email,
	})
}
