package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rashidmailru/kabobfood/internal/auth"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	service *auth.Service
}

// NewAuthHandler builds AuthHandler.
func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register wires the auth routes.
func (h *AuthHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/auth/telegram", h.telegramAuth)
}

type telegramAuthRequest struct {
	InitData string `json:"init_data" binding:"required"`
}

func (h *AuthHandler) telegramAuth(c *gin.Context) {
	var req telegramAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "init_data is required"})
		return
	}

	result, err := h.service.Authenticate(c.Request.Context(), req.InitData)
	if err != nil {
		switch err {
		case auth.ErrInvalidInitData, auth.ErrExpiredInitData, auth.ErrMissingUserPayload:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to authenticate"})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
