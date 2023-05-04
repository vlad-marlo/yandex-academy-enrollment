package main

import (
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/config"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller/http"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/middleware"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/logger"
	"go.uber.org/fx"
)

//	@title		Yandex Lavka
//	@version	1.0

func main() {
	fx.New(CreateApp()).Run()
}

// CreateApp prepares fx options to run server.
//
// This makes available to test is configuration correct.
func CreateApp() fx.Option {
	return fx.Options(
		fx.Provide(
			logger.New,
			fx.Annotate(http.New, fx.As(new(controller.Server))),
			fx.Annotate(config.NewRateLimiterConfig, fx.As(new(middleware.RateLimitConfig))),
			fx.Annotate(config.NewControllerConfig, fx.As(new(controller.Config))),
		),
		fx.Invoke(
			RunServer,
		),
		fx.NopLogger,
	)
}

// RunServer is helper function to configure server.
func RunServer(lc fx.Lifecycle, server controller.Server) {
	lc.Append(fx.Hook{
		OnStart: server.Start,
		OnStop:  server.Stop,
	})
}
