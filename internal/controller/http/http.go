package http

import (
	"context"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/vlad-marlo/yandex-academy-enrollment/docs"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	mw "github.com/vlad-marlo/yandex-academy-enrollment/internal/middleware"
	"go.uber.org/zap"
)

type Controller struct {
	engine  *echo.Echo
	log     *zap.Logger
	cfg     controller.Config
	srv     controller.Service
	rateCfg mw.RateLimitConfig
}

func New(
	logger *zap.Logger,
	cfg controller.Config,
	rateCfg mw.RateLimitConfig,
) (*Controller, error) {
	srv := &Controller{
		engine:  echo.New(),
		log:     logger,
		cfg:     cfg,
		rateCfg: rateCfg,
	}
	if logger == nil || cfg == nil || rateCfg == nil {
		return nil, ErrNilReference
	}
	srv.configure()
	logger.Info("successful initialized server")
	return srv, nil
}

func (srv *Controller) configureMW() {
	srv.engine.Use(
		mw.RateLimiter(srv.rateCfg),
	)
}

func (srv *Controller) configureRoutes() {
	srv.engine.GET("/swagger/*", echoSwagger.WrapHandler)
	srv.engine.GET("/ping", srv.HandlePing)
	couriers := srv.engine.Group("/couriers")
	couriers.GET("/:courier_id", srv.HandleGetCourier, mw.Paginator)
	couriers.POST("/", srv.HandleCreateCouriers)
}

func (srv *Controller) configure() {
	srv.configureMW()
	srv.configureRoutes()
}

func (srv *Controller) Start(context.Context) error {
	go func() {
		srv.log.Error("starting http server", zap.Error(srv.engine.Start(srv.cfg.BindAddr())))
	}()
	srv.log.Info("starting http server", zap.String("bind_addr", srv.cfg.BindAddr()))
	return nil
}

func (srv *Controller) Stop(ctx context.Context) error {
	srv.log.Info("stopping http server", zap.String("bind_addr", srv.cfg.BindAddr()))
	return srv.engine.Shutdown(ctx)
}
