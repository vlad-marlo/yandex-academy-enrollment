package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type RateLimiterConfig struct {
	bursts int `env:"RATE_LIMITER_BURSTS" envDefault:"1"`
}

// rateLimit is default limit of requests per second for each endpoint.
const rateLimit rate.Limit = 10

func NewRateLimiterConfig() (*RateLimiterConfig, error) {
	cfg := new(RateLimiterConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("env: parse: %w", err)
	}
	return cfg, nil
}

func (cfg *RateLimiterConfig) Burst() int {
	if cfg == nil {
		zap.L().Warn("unexpectedly got nil config object")
		return 0
	}
	return cfg.bursts
}

func (cfg *RateLimiterConfig) Limit() rate.Limit {
	return rateLimit
}
