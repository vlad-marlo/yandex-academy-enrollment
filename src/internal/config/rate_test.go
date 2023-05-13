package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewRateLimiterConfig(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		before := os.Getenv("RATE_LIMITER_BURSTS")
		defer assert.NoError(t, os.Setenv("RATE_LIMITER_BURSTS", before))
		assert.NoError(t, os.Setenv("RATE_LIMITER_BURSTS", "9"))
		cfg, err := NewRateLimiterConfig()
		assert.NoError(t, err)
		if assert.NotNil(t, cfg) {
			assert.Equal(t, 9, cfg.B)
			assert.Equal(t, rateLimit, cfg.Limit())
			assert.Equal(t, cfg.B, cfg.Burst())
		}
	})
	t.Run("bad burst", func(t *testing.T) {
		before := os.Getenv("RATE_LIMITER_BURSTS")
		defer assert.NoError(t, os.Setenv("RATE_LIMITER_BURSTS", before))
		assert.NoError(t, os.Setenv("RATE_LIMITER_BURSTS", "non-int"))
		cfg, err := NewRateLimiterConfig()
		assert.Error(t, err)
		if assert.Nil(t, cfg) {
			assert.Equal(t, rateLimit, cfg.Limit())
			assert.Equal(t, 0, cfg.Burst())
		}
	})
}
