package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/rashidmailru/kabobfood/internal/admin"
)

// AdminAuthHandler issues admin JWT tokens.
type AdminAuthHandler struct {
	authService *admin.AuthService
	jwtSecret   []byte
	jwtExpiry   time.Duration
}

// NewAdminAuthHandler builds handler.
func NewAdminAuthHandler(service *admin.AuthService, secret string, expiry time.Duration) *AdminAuthHandler {
	return &AdminAuthHandler{
		authService: service,
		jwtSecret:   []byte(secret),
		jwtExpiry:   expiry,
	}
}

// Register wires routes.
func (h *AdminAuthHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/admin/login", h.login)
}

type adminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AdminAuthHandler) login(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	user, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"role":  "admin",
		"iat":   now.Unix(),
		"exp":   now.Add(h.jwtExpiry).Unix(),
		"scope": "admin",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signed})
}
