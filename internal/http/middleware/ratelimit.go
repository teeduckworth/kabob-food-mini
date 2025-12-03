package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket per IP.
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	limit   int
	window  time.Duration
}

type bucket struct {
	tokens int
	reset  time.Time
}

// NewRateLimiter creates limiter with tokens per window.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{buckets: make(map[string]*bucket), limit: limit, window: window}
}

// Middleware returns Gin middleware enforcing rate limit.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	b, ok := rl.buckets[key]
	now := time.Now()
	if !ok || now.After(b.reset) {
		rl.buckets[key] = &bucket{tokens: rl.limit - 1, reset: now.Add(rl.window)}
		return true
	}
	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}
