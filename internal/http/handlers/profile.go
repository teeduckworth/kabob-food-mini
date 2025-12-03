package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rashidmailru/kabobfood/internal/http/middleware"
	"github.com/rashidmailru/kabobfood/internal/profile"
)

// ProfileHandler exposes profile endpoints.
type ProfileHandler struct {
	service *profile.Service
}

// NewProfileHandler creates handler.
func NewProfileHandler(service *profile.Service) *ProfileHandler {
	return &ProfileHandler{service: service}
}

// Register wires profile routes under authenticated group.
func (h *ProfileHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/profile", h.getProfile)
}

func (h *ProfileHandler) getProfile(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	data, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load profile"})
		return
	}

	c.JSON(http.StatusOK, data)
}
