package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	mw "github.com/vlad-marlo/yandex-academy-enrollment/internal/middleware"
	"go.uber.org/multierr"
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
	engine *echo.Echo,
	logger *zap.Logger,
	cfg controller.Config,
	rateCfg mw.RateLimitConfig,
) (*Controller, error) {
	srv := &Controller{
		engine:  engine,
		log:     logger,
		cfg:     cfg,
		rateCfg: rateCfg,
	}
	if engine == nil || logger == nil || cfg == nil || rateCfg == nil {
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
	srv.engine.GET("/ping", srv.HandlePing)
	couriers := srv.engine.Group("/couriers")
	couriers.GET("/:courier_id", srv.HandleGetCourier, mw.Paginator)
	couriers.POST("/", srv.HandleCreateCouriers)
}

func (srv *Controller) configure() {
	srv.configureMW()
	srv.configureRoutes()
}

func (srv *Controller) Start(context.Context) (err error) {
	srv.log.Info("starting http server", zap.String("bind_addr", srv.cfg.BindAddr()))
	go multierr.AppendInto(&err, srv.engine.Start(srv.cfg.BindAddr()))
	return
}

func (srv *Controller) Stop(ctx context.Context) error {
	srv.log.Info("stopping http server", zap.String("bind_addr", srv.cfg.BindAddr()))
	return srv.engine.Shutdown(ctx)
}
