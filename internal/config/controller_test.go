package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetControllerConfig(t *testing.T) {
	before := GetControllerConfig()
	assert.Equal(t, globalControllerConfig, before)
	globalControllerConfig = new(ControllerConfig)
	assert.NotEqual(t, before, GetControllerConfig())
}

func TestReplaceGlobalControllerConfig(t *testing.T) {
	before := GetControllerConfig()
	assert.Equal(t, globalControllerConfig, before)
	after := &ControllerConfig{Addr: "some random addr"}
	restore := ReplaceGlobalControllerConfig(after)
	assert.NotEqual(t, before, GetControllerConfig())
	restore()
	assert.Equal(t, before, GetControllerConfig())
}

func TestControllerConfig_BindAddr(t *testing.T) {
	tt := []struct {
		name string
		cfg  *ControllerConfig
		want string
	}{
		{"nil cfg", nil, defaultBindAddr},
		{"normal cfg", globalControllerConfig, globalControllerConfig.Addr},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.cfg.BindAddr())
		})
	}
}

func TestNewControllerConfig(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		before := os.Getenv("BIND_ADDR")
		defer assert.NoError(t, os.Setenv("BIND_ADDR", before))
		addr := "some addr"
		require.NoError(t, os.Setenv("BIND_ADDR", addr))
		cfg, err := NewControllerConfig()
		assert.NoError(t, err)
		if assert.NotNil(t, cfg) {
			assert.Equal(t, cfg.Addr, addr)
		}
	})
}
