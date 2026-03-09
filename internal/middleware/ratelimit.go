package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v5"
	"golang.org/x/time/rate"
)

// visitor holds the rate limiter and last seen time for a given IP
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// ipLimiter manages per-IP rate limiters
type ipLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	r        rate.Limit // requests per second
	b        int        // burst size
}

func newIPLimiter(r rate.Limit, b int) *ipLimiter {
	il := &ipLimiter{
		visitors: make(map[string]*visitor),
		r:        r,
		b:        b,
	}
	// Clean up stale visitors every 5 minutes
	go il.cleanupLoop()
	return il
}

func (il *ipLimiter) get(ip string) *rate.Limiter {
	il.mu.Lock()
	defer il.mu.Unlock()

	v, exists := il.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(il.r, il.b)
		il.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (il *ipLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		il.mu.Lock()
		for ip, v := range il.visitors {
			if time.Since(v.lastSeen) > 10*time.Minute {
				delete(il.visitors, ip)
			}
		}
		il.mu.Unlock()
	}
}

// RateLimit returns a middleware with the given rate (per second) and burst.
// When the limit is exceeded it returns 429 with a plain message.
func RateLimit(r rate.Limit, burst int) echo.MiddlewareFunc {
	limiter := newIPLimiter(r, burst)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ip := c.RealIP()
			if !limiter.get(ip).Allow() {
				return c.String(http.StatusTooManyRequests, "Too many requests. Please slow down.")
			}
			return next(c)
		}
	}
}
