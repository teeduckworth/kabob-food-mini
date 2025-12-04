package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rashidmailru/kabobfood/internal/auth"
)

// AuthHandler handles Telegram WebApp authentication requests.
type AuthHandler struct {
	service *auth.Service
}

// NewAuthHandler constructs AuthHandler.
func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register wires auth routes.
func (h *AuthHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/auth/telegram", h.authenticateTelegram)
}

type telegramAuthRequest struct {
	InitData string `json:"init_data" binding:"required"`
}

func (h *AuthHandler) authenticateTelegram(c *gin.Context) {
	var req telegramAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "init_data is required"})
		return
	}

	result, err := h.service.Authenticate(c.Request.Context(), req.InitData)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidInitData) || errors.Is(err, auth.ErrExpiredInitData) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate"})
		return
	}

	c.JSON(http.StatusOK, result)
}
