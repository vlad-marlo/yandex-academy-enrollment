package main

import (
	"github.com/labstack/echo"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller/http"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/logger"
	"go.uber.org/fx"
)

func main() {
	fx.New(CreateApp()).Run()
}

func CreateApp() fx.Option {
	return fx.Options(
		fx.Provide(
			logger.New,
			echo.New,
			fx.Annotate(http.New, fx.As(new(controller.Server))),
		),
		fx.Provide(RunServer),
	)
}

func RunServer(lc fx.Lifecycle, server controller.Server) {
	lc.Append(fx.Hook{
		OnStart: server.Start,
		OnStop:  server.Stop,
	})
}
