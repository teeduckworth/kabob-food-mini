package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rashidmailru/kabobfood/internal/menu"
)

// MenuHandler exposes menu + regions endpoints for Mini App.
type MenuHandler struct {
	service *menu.Service
}

// NewMenuHandler creates handler.
func NewMenuHandler(service *menu.Service) *MenuHandler {
	return &MenuHandler{service: service}
}

// Register wires routes.
func (h *MenuHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/menu", h.getMenu)
	rg.GET("/regions", h.getRegions)
}

func (h *MenuHandler) getMenu(c *gin.Context) {
	resp, err := h.service.GetMenu(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load menu"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MenuHandler) getRegions(c *gin.Context) {
	resp, err := h.service.GetRegions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load regions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"regions": resp})
}
