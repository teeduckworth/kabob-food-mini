package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rashidmailru/kabobfood/internal/orders"
)

// AdminOrdersHandler exposes operator endpoints for orders.
type AdminOrdersHandler struct {
	service *orders.AdminService
}

func NewAdminOrdersHandler(service *orders.AdminService) *AdminOrdersHandler {
	return &AdminOrdersHandler{service: service}
}

func (h *AdminOrdersHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/admin/orders", h.list)
	rg.PUT("/admin/orders/:id/status", h.updateStatus)
}

func (h *AdminOrdersHandler) list(c *gin.Context) {
	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}
	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsed
		}
	}
	params := orders.AdminListParams{
		Status: c.Query("status"),
		Limit:  limit,
		Offset: offset,
	}
	if fromStr := c.Query("from"); fromStr != "" {
		if ts, err := time.Parse(time.RFC3339, fromStr); err == nil {
			params.From = &ts
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if ts, err := time.Parse(time.RFC3339, toStr); err == nil {
			params.To = &ts
		}
	}
	ordersList, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orders": ordersList})
}

func (h *AdminOrdersHandler) updateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	var req orders.UpdateStatusInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	order, err := h.service.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}
