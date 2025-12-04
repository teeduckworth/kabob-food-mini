package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/rashidmailru/kabobfood/internal/auth"
)

// BotHandler handles registration requests originating from the Telegram bot.
type BotHandler struct {
	service *auth.Service
}

// NewBotHandler builds BotHandler.
func NewBotHandler(service *auth.Service) *BotHandler {
	return &BotHandler{service: service}
}

// Register wires bot-specific routes.
func (h *BotHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/bot/register", h.register)
}

type botRegisterRequest struct {
	TelegramID int64               `json:"telegram_id" binding:"required"`
	Phone      string              `json:"phone" binding:"required"`
	FirstName  string              `json:"first_name"`
	LastName   string              `json:"last_name"`
	Name       string              `json:"name"`
	Location   *botLocationPayload `json:"location" binding:"required"`
}

type botLocationPayload struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (h *BotHandler) register(c *gin.Context) {
	var req botRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	firstName := strings.TrimSpace(req.FirstName)
	if firstName == "" {
		firstName = strings.TrimSpace(req.Name)
	}
	lastName := strings.TrimSpace(req.LastName)
	phone := strings.TrimSpace(req.Phone)
	if firstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "first name is required"})
		return
	}
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone is required"})
		return
	}
	if req.Location == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "location is required"})
		return
	}

	lat := req.Location.Latitude
	lon := req.Location.Longitude

	result, err := h.service.RegisterBotUser(c.Request.Context(), auth.BotRegisterInput{
		TelegramID: req.TelegramID,
		FirstName:  firstName,
		LastName:   lastName,
		Phone:      phone,
		Latitude:   &lat,
		Longitude:  &lon,
	})
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRegisterInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		return
	}

	c.JSON(http.StatusOK, result)
}
