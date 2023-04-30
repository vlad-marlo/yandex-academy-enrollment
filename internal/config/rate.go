package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimiterConfig implements middleware.RateLimiterConfig type.
type RateLimiterConfig struct {
	B int `env:"RATE_LIMITER_BURSTS" envDefault:"10"`
}

// rateLimit is default limit of requests per second for each endpoint.
const rateLimit rate.Limit = 10

// NewRateLimiterConfig configures rate limiter.
func NewRateLimiterConfig() (*RateLimiterConfig, error) {
	cfg := new(RateLimiterConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("env: parse: %w", err)
	}
	return cfg, nil
}

// Burst returns burst for rate limiter.
func (cfg *RateLimiterConfig) Burst() int {
	if cfg == nil {
		zap.L().Warn("unexpectedly got nil config object")
		return 0
	}
	return cfg.B
}

// Limit returns limit for rate limiter.
func (*RateLimiterConfig) Limit() rate.Limit {
	return rateLimit
}
