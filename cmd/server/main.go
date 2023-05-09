package main

import (
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/config"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller/http"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/middleware"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/logger"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/pgx"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/pgx/client"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/pgx/migrator"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/service/example"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/service/production"
	pgxStore "github.com/vlad-marlo/yandex-academy-enrollment/internal/store/pgx"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
			fx.Annotate(config.NewPgConfig, fx.As(new(client.Config))),
			fx.Annotate(client.New, fx.As(new(pgx.Client))),
			fx.Annotate(pgxStore.New, fx.As(new(production.Store))),
			fx.Annotate(example.New, fx.As(new(controller.Service))),
		),
		fx.Invoke(
			RunServer,
			Migrate,
		),
		fx.NopLogger,
	)
}

func Migrate(cli pgx.Client) error {
	migrations, err := migrator.Migrate(cli)
	cli.L().Info("migrated database", zap.Int("migrations_applied", migrations))
	return err
}

// RunServer is helper function to configure server.
func RunServer(lc fx.Lifecycle, server controller.Server) {
	lc.Append(fx.Hook{
		OnStart: server.Start,
		OnStop:  server.Stop,
	})
}
