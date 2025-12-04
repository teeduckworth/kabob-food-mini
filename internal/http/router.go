package http

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/rashidmailru/kabobfood/internal/http/handlers"
	"github.com/rashidmailru/kabobfood/internal/metrics"
)

// RouteRegister defines interface implemented by handlers with Register method.
type RouteRegister interface {
	Register(*gin.RouterGroup)
}

// RouterParams groups construction dependencies for the HTTP router.
type RouterParams struct {
	Logger            *zap.Logger
	AppEnv            string
	CORSOrigins       []string
	HealthHandler     *handlers.HealthHandler
	BotHandler        *handlers.BotHandler
	MenuHandler       *handlers.MenuHandler
	AdminAuthHandler  *handlers.AdminAuthHandler
	AuthMiddleware    gin.HandlerFunc
	ProtectedHandlers []RouteRegister
	AdminMiddleware   gin.HandlerFunc
	AdminHandlers     []RouteRegister
	Metrics           *metrics.Metrics
}

// NewRouter configures Gin engine with logging, recovery, and base routes.
func NewRouter(params RouterParams) *gin.Engine {
	if params.AppEnv == "prod" || params.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	corsCfg := buildCORSConfig(params.CORSOrigins)

	router.Use(ginzap.Ginzap(params.Logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(params.Logger, true))
	router.Use(cors.New(corsCfg))
	if params.Metrics != nil {
		router.Use(params.Metrics.Middleware())
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	root := router.Group("")
	params.HealthHandler.Register(root)
	if params.BotHandler != nil {
		params.BotHandler.Register(root)
	}
	if params.AdminAuthHandler != nil {
		params.AdminAuthHandler.Register(root)
	}
	if params.MenuHandler != nil {
		params.MenuHandler.Register(root)
	}

	if params.AuthMiddleware != nil && len(params.ProtectedHandlers) > 0 {
		authGroup := root.Group("")
		authGroup.Use(params.AuthMiddleware)
		for _, h := range params.ProtectedHandlers {
			if h != nil {
				h.Register(authGroup)
			}
		}
	}

	if params.AdminMiddleware != nil && len(params.AdminHandlers) > 0 {
		adminGroup := root.Group("")
		adminGroup.Use(params.AdminMiddleware)
		for _, h := range params.AdminHandlers {
			if h != nil {
				h.Register(adminGroup)
			}
		}
	}

	return router
}

func buildCORSConfig(origins []string) cors.Config {
	cleanOrigins := normalizeOrigins(origins)
	cfg := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}
	if len(cleanOrigins) == 0 {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = cleanOrigins
	}
	return cfg
}

func normalizeOrigins(origins []string) []string {
	if len(origins) == 0 {
		return nil
	}
	clean := make([]string, 0, len(origins))
	seen := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		trimmed := strings.ToLower(strings.TrimSpace(origin))
		trimmed = strings.TrimRight(trimmed, "/")
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		clean = append(clean, trimmed)
	}
	if len(clean) == 0 {
		return nil
	}
	return clean
}
