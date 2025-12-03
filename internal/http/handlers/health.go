package handlers

import "github.com/gin-gonic/gin"

// HealthHandler reports service health and metadata.
type HealthHandler struct {
	Version string
}

// NewHealthHandler constructs a HealthHandler.
func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{Version: version}
}

// Register wires handler routes under the given router group.
func (h *HealthHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/healthz", h.health)
	rg.GET("/version", h.version)
}

func (h *HealthHandler) health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func (h *HealthHandler) version(c *gin.Context) {
	c.JSON(200, gin.H{"version": h.Version})
}
