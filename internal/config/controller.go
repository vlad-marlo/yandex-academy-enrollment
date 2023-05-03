package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"go.uber.org/zap"
	"sync"
)

const (
	defaultBindAddr = "localhost:8080"
)

var (
	globalControllerConfig = &ControllerConfig{
		Addr: defaultBindAddr,
	}
	globalControllerMu sync.RWMutex
)

// ControllerConfig configures web server.
type ControllerConfig struct {
	// Addr specifies address on which server will be running at.
	Addr string `env:"BIND_ADDR" envDefault:"localhost:8080"`
}

// NewControllerConfig initializes controller config and returns it to user.
func NewControllerConfig() (*ControllerConfig, error) {
	cfg := new(ControllerConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("env: parse: %w", err)
	}
	return cfg, nil
}

// GetControllerConfig returns global config.
func GetControllerConfig() *ControllerConfig {
	globalControllerMu.RLock()
	cfg := globalControllerConfig
	globalControllerMu.RUnlock()
	return cfg
}

// ReplaceGlobalControllerConfig replaces global config with provided.
func ReplaceGlobalControllerConfig(cfg *ControllerConfig) func() {
	globalControllerMu.Lock()
	before := globalControllerConfig
	globalControllerConfig = cfg
	globalControllerMu.Unlock()
	return func() {
		ReplaceGlobalControllerConfig(before)
	}
}

// BindAddr returns address on which server will be running at.
func (cfg *ControllerConfig) BindAddr() string {
	if cfg == nil {
		zap.L().Warn("unexpectedly got nil pointer receiver config")
		return defaultBindAddr
	}
	return cfg.Addr
}
