package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

/**
 * visitor tracks a single client's token-bucket state.
 */
type visitor struct {
	tokens   float64
	lastSeen time.Time
}

/**
 * rateLimiter is an in-memory, per-IP token-bucket rate limiter.
 * No external dependency and no shared store — fine for a single instance.
 */
type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     float64 // tokens reabastecidos por segundo (vazão sustentada)
	burst    float64 // capacidade máxima do balde (pico)
}

func newRateLimiter(rate, burst float64) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}
	go rl.cleanup()
	return rl
}

/**
 * cleanup evicts idle visitors every 10 minutes to bound memory.
 */
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 15*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

/**
 * allow consumes one token for the given IP, refilling based on elapsed time.
 * Returns false when the bucket is empty (request should be rejected).
 */
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, ok := rl.visitors[ip]
	if !ok {
		rl.visitors[ip] = &visitor{tokens: rl.burst - 1, lastSeen: now}
		return true
	}

	// Reabastece conforme o tempo decorrido, limitado ao burst.
	v.tokens += now.Sub(v.lastSeen).Seconds() * rl.rate
	if v.tokens > rl.burst {
		v.tokens = rl.burst
	}
	v.lastSeen = now

	if v.tokens < 1 {
		return false
	}
	v.tokens--
	return true
}

/**
 * RateLimit returns a middleware that limits requests per client IP.
 *
 * @param ratePerSec - vazão sustentada (tokens por segundo)
 * @param burst - número máximo de requisições em pico
 */
func RateLimit(ratePerSec, burst float64) gin.HandlerFunc {
	rl := newRateLimiter(ratePerSec, burst)
	return func(c *gin.Context) {
		if !rl.allow(c.ClientIP()) {
			c.Header("Retry-After", "60")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Muitas requisições. Tente novamente em instantes.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
