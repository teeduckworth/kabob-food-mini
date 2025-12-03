package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rashidmailru/kabobfood/internal/admin"
	"github.com/rashidmailru/kabobfood/internal/menu"
)

// AdminMenuHandler exposes admin CRUD endpoints for menu entities.
type AdminMenuHandler struct {
	service *admin.MenuService
}

func NewAdminMenuHandler(service *admin.MenuService) *AdminMenuHandler {
	return &AdminMenuHandler{service: service}
}

func (h *AdminMenuHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/admin/categories", h.createCategory)
	rg.PUT("/admin/categories/:id", h.updateCategory)
	rg.DELETE("/admin/categories/:id", h.deleteCategory)

	rg.POST("/admin/products", h.createProduct)
	rg.PUT("/admin/products/:id", h.updateProduct)
	rg.DELETE("/admin/products/:id", h.deleteProduct)
}

func (h *AdminMenuHandler) createCategory(c *gin.Context) {
	var req menu.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	created, err := h.service.CreateCategory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *AdminMenuHandler) updateCategory(c *gin.Context) {
	var req menu.Category
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
	updated, err := h.service.UpdateCategory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *AdminMenuHandler) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AdminMenuHandler) createProduct(c *gin.Context) {
	var req menu.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	created, err := h.service.CreateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *AdminMenuHandler) updateProduct(c *gin.Context) {
	var req menu.Product
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
	updated, err := h.service.UpdateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *AdminMenuHandler) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
