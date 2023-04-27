package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	mw "github.com/vlad-marlo/yandex-academy-enrollment/internal/middleware"
	"go.uber.org/zap"
)

type Controller struct {
	engine           *echo.Echo
	log              *zap.Logger
	cfg              controller.Config
	rateCfg          mw.RateLimitConfig
	routesConfigured bool
}

func New(
	engine *echo.Echo,
	logger *zap.Logger,
	cfg controller.Config,
	rateCfg mw.RateLimitConfig,
) (*Controller, error) {
	srv := &Controller{
		engine: engine,
		log:    logger,
		cfg:    cfg,
	}
	switch nil {
	case engine, logger, cfg, rateCfg:
		return nil, ErrNilReference
	default:
	}
	return srv, nil
}

func (srv *Controller) configureMW() {
	if srv.routesConfigured {
		srv.engine.Use()
	}
}

func (srv *Controller) Start(context.Context) error {
	return srv.engine.Start(srv.cfg.BindAddr())
}

func (srv *Controller) Stop(ctx context.Context) error {
	return srv.engine.Shutdown(ctx)
}
