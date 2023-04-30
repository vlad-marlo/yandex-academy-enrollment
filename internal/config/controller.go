package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"go.uber.org/zap"
)

type ControllerConfig struct {
	Addr string `env:"BIND_ADDR" envDefault:"localhost:8080"`
}

func NewControllerConfig() (*ControllerConfig, error) {
	cfg := new(ControllerConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("env: parse: %w", err)
	}
	return cfg, nil
}

func (cfg *ControllerConfig) BindAddr() string {
	if cfg == nil {
		zap.L().Warn("unexpectedly got nil pointer receiver config")
		return ""
	}
	return cfg.Addr
}
