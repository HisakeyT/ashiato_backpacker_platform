package handler

import (
	"net/http"

	"github.com/HisakeyT/backpacker-platform/internal/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userUsecase *usecase.UserUsecase
}

func NewHandler(userUsecase *usecase.UserUsecase) *Handler {
	return &Handler{
		userUsecase: userUsecase,
	}
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

	u, err := h.userUsecase.Register(input)
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

	result, err := h.userUsecase.Login(input)
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

	user, err := h.userUsecase.Me(userID.(uint))
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

	user, err := h.userUsecase.UpdateMe(userID.(uint), input)
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
