package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/rashidmailru/kabobfood/internal/config"
)

// NewRedis builds a redis client and validates the connection with a ping.
func NewRedis(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	var (
		opts *redis.Options
		err  error
	)

	opts, err = redis.ParseURL(cfg.URL)
	if err != nil {
		opts = &redis.Options{Addr: cfg.URL}
	}

	if cfg.Password != "" {
		opts.Password = cfg.Password
	}

	client := redis.NewClient(opts)

	timeout := cfg.DialTimeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
