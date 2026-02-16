package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type windowEntry struct {
	count       atomic.Int64
	windowStart atomic.Int64
}

type RateLimitMiddleware struct {
	entries   sync.Map
	limit     int64
	windowNs  int64
	stopClean chan struct{}
}

func NewRateLimitMiddleware(maxRequests int, window time.Duration) *RateLimitMiddleware {
	rl := &RateLimitMiddleware{
		limit:     int64(maxRequests),
		windowNs:  window.Nanoseconds(),
		stopClean: make(chan struct{}),
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimitMiddleware) Stop() {
	close(rl.stopClean)
}

func (rl *RateLimitMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if projectID := c.GetHeader("X-Project-Id"); projectID != "" {
			key = key + ":" + projectID
		}

		now := time.Now().UnixNano()

		val, _ := rl.entries.LoadOrStore(key, &windowEntry{})
		entry := val.(*windowEntry)

		windowStart := entry.windowStart.Load()
		if windowStart == 0 || now-windowStart >= rl.windowNs {
			entry.windowStart.Store(now)
			entry.count.Store(1)

			rl.setHeaders(c, 1, now)
			c.Next()
			return
		}

		count := entry.count.Add(1)

		if count > rl.limit {
			remaining := windowStart + rl.windowNs - now
			retryAfter := time.Duration(remaining).Seconds()
			if retryAfter < 1 {
				retryAfter = 1
			}

			rl.setHeaders(c, count, windowStart)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate limit exceeded",
				"retry_after": fmt.Sprintf("%.0fs", retryAfter),
			})
			return
		}

		rl.setHeaders(c, count, windowStart)
		c.Next()
	}
}

func (rl *RateLimitMiddleware) setHeaders(c *gin.Context, count int64, windowStart int64) {
	remaining := rl.limit - count
	if remaining < 0 {
		remaining = 0
	}

	resetAt := time.Unix(0, windowStart+rl.windowNs).Unix()

	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
	c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", resetAt))
}

func (rl *RateLimitMiddleware) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-rl.stopClean:
			return
		case <-ticker.C:
			now := time.Now().UnixNano()
			rl.entries.Range(func(key, value any) bool {
				entry := value.(*windowEntry)
				if now-entry.windowStart.Load() >= rl.windowNs {
					rl.entries.Delete(key)
				}
				return true
			})
		}
	}
}
