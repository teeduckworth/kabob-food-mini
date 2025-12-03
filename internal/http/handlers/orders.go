package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rashidmailru/kabobfood/internal/http/middleware"
	"github.com/rashidmailru/kabobfood/internal/orders"
)

// OrdersHandler handles order-related endpoints for users.
type OrdersHandler struct {
	service *orders.Service
}

// NewOrdersHandler constructs handler.
func NewOrdersHandler(service *orders.Service) *OrdersHandler {
	return &OrdersHandler{service: service}
}

// Register wires order routes (auth required).
func (h *OrdersHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/orders", h.create)
	rg.GET("/orders", h.list)
	rg.GET("/orders/:id", h.get)
}

func (h *OrdersHandler) create(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req orders.CreateOrderInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	order, err := h.service.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrdersHandler) list(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ordersList, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": ordersList})
}

func (h *OrdersHandler) get(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.service.Get(c.Request.Context(), userID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}
