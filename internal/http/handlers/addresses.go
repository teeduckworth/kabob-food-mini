package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rashidmailru/kabobfood/internal/addresses"
	"github.com/rashidmailru/kabobfood/internal/http/middleware"
)

// AddressesHandler manages CRUD endpoints for user addresses.
type AddressesHandler struct {
	service *addresses.Service
}

// NewAddressesHandler constructs handler.
func NewAddressesHandler(service *addresses.Service) *AddressesHandler {
	return &AddressesHandler{service: service}
}

// Register wires address routes.
func (h *AddressesHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/addresses", h.list)
	rg.POST("/addresses", h.create)
	rg.PUT("/addresses/:id", h.update)
	rg.DELETE("/addresses/:id", h.delete)
}

func (h *AddressesHandler) list(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	items, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load addresses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"addresses": items})
}

type addressRequest struct {
	RegionID  int64  `json:"region_id" binding:"required"`
	Street    string `json:"street" binding:"required"`
	House     string `json:"house" binding:"required"`
	Entrance  string `json:"entrance"`
	Flat      string `json:"flat"`
	Comment   string `json:"comment"`
	IsDefault bool   `json:"is_default"`
}

func (h *AddressesHandler) create(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req addressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	addr, err := h.service.Create(c.Request.Context(), addresses.CreateInput{
		UserID:    userID,
		RegionID:  req.RegionID,
		Street:    req.Street,
		House:     req.House,
		Entrance:  req.Entrance,
		Flat:      req.Flat,
		Comment:   req.Comment,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, addr)
}

func (h *AddressesHandler) update(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	addressID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address id"})
		return
	}

	var req addressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	addr, err := h.service.Update(c.Request.Context(), addresses.UpdateInput{
		ID:        addressID,
		UserID:    userID,
		RegionID:  req.RegionID,
		Street:    req.Street,
		House:     req.House,
		Entrance:  req.Entrance,
		Flat:      req.Flat,
		Comment:   req.Comment,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update address"})
		return
	}

	c.JSON(http.StatusOK, addr)
}

func (h *AddressesHandler) delete(c *gin.Context) {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	addressID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), addressID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete address"})
		return
	}

	c.Status(http.StatusNoContent)
}
