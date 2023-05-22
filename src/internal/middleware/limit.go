package middleware

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

// RateLimitConfig is config of rate limiter mw object.
type RateLimitConfig interface {
	// Limit returns the maximum overall event rate.
	Limit() rate.Limit
	// Burst returns the maximum burst size. Burst is the maximum number of tokens
	// that can be consumed in a single call to Allow, Reserve, or Wait, so higher
	// Burst values allow more events to happen at once.
	// A zero Burst allows no events, unless limit == Inf.
	Burst() int
}

var (
	// limiters are limiters for each endpoint with access by url path.
	limiters = make(map[string]*rate.Limiter)
	// rlmu is global rw mutex of RateLimiter mw.
	rlmu sync.RWMutex
)

// RateLimiter add limit of RPS for every single route.
func RateLimiter(cfg RateLimitConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rlmu.RLock()
			limiter, ok := limiters[c.Request().URL.Path]
			rlmu.RUnlock()
			if !ok {
				limiter = rate.NewLimiter(cfg.Limit(), cfg.Burst())
				rlmu.Lock()
				limiters[c.Request().URL.Path] = limiter
				rlmu.Unlock()
			}
			if limiter.Allow() {
				return next(c)
			}
			return c.String(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
		}
	}
}
