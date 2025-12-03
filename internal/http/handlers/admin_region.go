package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rashidmailru/kabobfood/internal/admin"
	"github.com/rashidmailru/kabobfood/internal/regions"
)

// AdminRegionHandler manages region CRUD.
type AdminRegionHandler struct {
	service *admin.RegionService
}

func NewAdminRegionHandler(service *admin.RegionService) *AdminRegionHandler {
	return &AdminRegionHandler{service: service}
}

func (h *AdminRegionHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/admin/regions", h.createRegion)
	rg.PUT("/admin/regions/:id", h.updateRegion)
	rg.DELETE("/admin/regions/:id", h.deleteRegion)
}

func (h *AdminRegionHandler) createRegion(c *gin.Context) {
	var req regions.Region
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	created, err := h.service.CreateRegion(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *AdminRegionHandler) updateRegion(c *gin.Context) {
	var req regions.Region
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	req.ID = id
	updated, err := h.service.UpdateRegion(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *AdminRegionHandler) deleteRegion(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteRegion(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
