package http

import (
	"time"

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
	HealthHandler     *handlers.HealthHandler
	AuthHandler       *handlers.AuthHandler
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

	router.Use(ginzap.Ginzap(params.Logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(params.Logger, true))
	if params.Metrics != nil {
		router.Use(params.Metrics.Middleware())
		router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	root := router.Group("")
	params.HealthHandler.Register(root)
	if params.AuthHandler != nil {
		params.AuthHandler.Register(root)
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
