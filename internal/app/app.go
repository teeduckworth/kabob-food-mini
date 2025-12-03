package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/rashidmailru/kabobfood/internal/addresses"
	"github.com/rashidmailru/kabobfood/internal/admin"
	"github.com/rashidmailru/kabobfood/internal/auth"
	cachepkg "github.com/rashidmailru/kabobfood/internal/cache"
	"github.com/rashidmailru/kabobfood/internal/config"
	"github.com/rashidmailru/kabobfood/internal/db"
	kabobhttp "github.com/rashidmailru/kabobfood/internal/http"
	"github.com/rashidmailru/kabobfood/internal/http/handlers"
	"github.com/rashidmailru/kabobfood/internal/http/middleware"
	"github.com/rashidmailru/kabobfood/internal/menu"
	"github.com/rashidmailru/kabobfood/internal/metrics"
	"github.com/rashidmailru/kabobfood/internal/notifications"
	"github.com/rashidmailru/kabobfood/internal/orders"
	"github.com/rashidmailru/kabobfood/internal/products"
	"github.com/rashidmailru/kabobfood/internal/profile"
	"github.com/rashidmailru/kabobfood/internal/regions"
	"github.com/rashidmailru/kabobfood/internal/server"
	"github.com/rashidmailru/kabobfood/internal/users"
)

// Version indicates the application build identifier.
const Version = "0.1.0"

// App wires the different layers together.
type App struct {
	cfg    *config.Config
	log    *zap.Logger
	server *server.Server
	dbPool *pgxpool.Pool
	cache  *redis.Client
}

// New constructs the application.
func New(cfg *config.Config, log *zap.Logger) (*App, error) {
	ctx := context.Background()
	pool, err := db.NewPostgres(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	redisClient, err := cachepkg.NewRedis(ctx, cfg.Redis)
	if err != nil {
		pool.Close()
		return nil, err
	}

	metricsCollector := metrics.New()

	userRepo := users.NewRepository(pool)
	addressRepo := addresses.NewRepository(pool)
	addressService := addresses.NewService(addressRepo)
	menuRepo := menu.NewRepository(pool)
	productsRepo := products.NewRepository(pool)
	regionRepo := regions.NewRepository(pool)
	adminRepo := admin.NewRepository(pool)

	authService, err := auth.NewService(auth.Config{
		UserRepo:    userRepo,
		BotToken:    cfg.Telegram.BotToken,
		JWTSecret:   cfg.JWT.Secret,
		JWTExpiry:   cfg.JWT.Expiration,
		InitDataTTL: cfg.Auth.TelegramInitTTL,
	})
	if err != nil {
		pool.Close()
		redisClient.Close()
		return nil, err
	}

	menuService := menu.NewService(menu.ServiceConfig{
		MenuRepo:   menuRepo,
		RegionRepo: regionRepo,
		Cache:      redisClient,
		MenuTTL:    cfg.Cache.MenuTTL,
		RegionsTTL: cfg.Cache.RegionsTTL,
	})
	adminMenuService := admin.NewMenuService(menuRepo, redisClient)
	adminRegionService := admin.NewRegionService(regionRepo, redisClient)

	profileService := profile.NewService(userRepo, addressService)
	ordersRepo := orders.NewRepository(pool)
	notifier := notifications.NewTelegramNotifier(notifications.TelegramConfig{
		BotToken:    cfg.Telegram.BotToken,
		AdminChatID: cfg.Telegram.AdminChatID,
	})
	ordersService := orders.NewService(ordersRepo, productsRepo, addressRepo, regionRepo, userRepo, notifier, metricsCollector)
	adminOrdersService := orders.NewAdminService(ordersRepo, userRepo, notifier)
	adminAuthService, err := admin.NewAuthService(admin.AuthConfig{Repo: adminRepo, JWTSecret: cfg.JWT.Secret})
	if err != nil {
		pool.Close()
		redisClient.Close()
		return nil, err
	}
	if err := adminAuthService.EnsureDefaultAdmin(ctx, cfg.Admin.DefaultUsername, cfg.Admin.DefaultPassword); err != nil {
		pool.Close()
		redisClient.Close()
		return nil, err
	}
	rateLimiterUsers := middleware.NewRateLimiter(cfg.RateLimit.UserLimit, cfg.RateLimit.Window)
	rateLimiterAdmins := middleware.NewRateLimiter(cfg.RateLimit.AdminLimit, cfg.RateLimit.Window)

	profileHandler := handlers.NewProfileHandler(profileService)
	addressesHandler := handlers.NewAddressesHandler(addressService)
	ordersHandler := handlers.NewOrdersHandler(ordersService)
	protectedHandlers := []kabobhttp.RouteRegister{profileHandler, addressesHandler, ordersHandler}
	jwtMiddleware := middleware.JWTAuth(cfg.JWT.Secret)
	adminAuthHandler := handlers.NewAdminAuthHandler(adminAuthService, cfg.JWT.Secret, cfg.Admin.JWTExpiration)
	adminMenuHandler := handlers.NewAdminMenuHandler(adminMenuService)
	adminRegionHandler := handlers.NewAdminRegionHandler(adminRegionService)
	adminOrdersHandler := handlers.NewAdminOrdersHandler(adminOrdersService)
	adminHandlers := []kabobhttp.RouteRegister{adminMenuHandler, adminRegionHandler, adminOrdersHandler}
	adminMiddleware := middleware.AdminJWT(cfg.JWT.Secret)

	healthHandler := handlers.NewHealthHandler(Version)
	authHandler := handlers.NewAuthHandler(authService)
	menuHandler := handlers.NewMenuHandler(menuService)

	router := kabobhttp.NewRouter(kabobhttp.RouterParams{
		Logger:            log,
		AppEnv:            cfg.AppEnv,
		HealthHandler:     healthHandler,
		AuthHandler:       authHandler,
		MenuHandler:       menuHandler,
		AdminAuthHandler:  adminAuthHandler,
		AuthMiddleware:    middleware.Chain(rateLimiterUsers.Middleware(), jwtMiddleware),
		ProtectedHandlers: protectedHandlers,
		AdminMiddleware:   middleware.Chain(rateLimiterAdmins.Middleware(), adminMiddleware),
		AdminHandlers:     adminHandlers,
		Metrics:           metricsCollector,
	})

	srv := server.New(cfg, router, log)

	return &App{cfg: cfg, log: log, server: srv, dbPool: pool, cache: redisClient}, nil
}

// Run starts handling HTTP traffic.
func (a *App) Run() error {
	return a.server.Run()
}

// Shutdown gracefully stops the server.
func (a *App) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	if a.dbPool != nil {
		a.dbPool.Close()
	}
	if a.cache != nil {
		_ = a.cache.Close()
	}
	return nil
}
