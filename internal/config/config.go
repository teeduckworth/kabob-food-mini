package config

import (
	"time"

	"github.com/caarlos0/env/v10"
)

// Config aggregates application configuration sourced from environment variables.
type Config struct {
	AppEnv          string          `env:"APP_ENV" envDefault:"development"`
	HTTP            HTTPConfig      `envPrefix:"HTTP_"`
	Database        DatabaseConfig  `envPrefix:"DB_"`
	Redis           RedisConfig     `envPrefix:"REDIS_"`
	JWT             JWTConfig       `envPrefix:"JWT_"`
	Telegram        TelegramConfig  `envPrefix:"TELEGRAM_"`
	Auth            AuthConfig      `envPrefix:"AUTH_"`
	Cache           CacheConfig     `envPrefix:"CACHE_"`
	Admin           AdminConfig     `envPrefix:"ADMIN_"`
	RateLimit       RateLimitConfig `envPrefix:"RATE_"`
	Sentry          SentryConfig    `envPrefix:"SENTRY_"`
	ShutdownTimeout time.Duration   `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

// HTTPConfig contains HTTP server settings.
type HTTPConfig struct {
	Host         string        `env:"HOST" envDefault:"0.0.0.0"`
	Port         int           `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
	CORSOrigins  []string      `env:"CORS_ORIGINS" envSeparator:","`
}

// DatabaseConfig declares PostgreSQL connection options.
type DatabaseConfig struct {
	URL             string        `env:"URL" envDefault:"postgres://postgres:postgres@postgres:5432/kabobfood?sslmode=disable"`
	MaxConns        int           `env:"MAX_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" envDefault:"1h"`
}

// RedisConfig declares redis connection options.
type RedisConfig struct {
	URL         string        `env:"URL"`
	Password    string        `env:"PASSWORD"`
	DialTimeout time.Duration `env:"DIAL_TIMEOUT" envDefault:"5s"`
}

// JWTConfig contains signing settings.
type JWTConfig struct {
	Secret     string        `env:"SECRET" envDefault:"supersecret"`
	Expiration time.Duration `env:"EXPIRATION" envDefault:"24h"`
}

// TelegramConfig defines credentials for Bot API integration.
type TelegramConfig struct {
	BotToken    string `env:"BOT_TOKEN"`
	AdminChatID string `env:"ADMIN_CHAT_ID"`
}

// AuthConfig includes various authentication settings.
type AuthConfig struct {
	TelegramInitTTL time.Duration `env:"TELEGRAM_INIT_TTL" envDefault:"1h"`
}

// CacheConfig defines TTLs for cached payloads.
type CacheConfig struct {
	MenuTTL    time.Duration `env:"MENU_TTL" envDefault:"30s"`
	RegionsTTL time.Duration `env:"REGIONS_TTL" envDefault:"30s"`
}

// AdminConfig defines bootstrap admin credentials and JWT ttl.
type AdminConfig struct {
	DefaultUsername string        `env:"DEFAULT_USERNAME" envDefault:"admin"`
	DefaultPassword string        `env:"DEFAULT_PASSWORD" envDefault:"admin123"`
	JWTExpiration   time.Duration `env:"JWT_EXPIRATION" envDefault:"24h"`
}

// SentryConfig stores sentry DSN.
type SentryConfig struct {
	DSN string `env:"DSN"`
}

// RateLimitConfig defines token bucket settings.
type RateLimitConfig struct {
	UserLimit  int           `env:"USER_LIMIT" envDefault:"60"`
	AdminLimit int           `env:"ADMIN_LIMIT" envDefault:"120"`
	Window     time.Duration `env:"WINDOW" envDefault:"1m"`
}

// Load parses environment variables into Config.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// MustLoad is a helper that panics if configuration cannot be loaded.
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
