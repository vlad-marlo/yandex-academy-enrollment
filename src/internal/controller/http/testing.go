package http

import (
	"github.com/labstack/echo/v4"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"testing"
)

const bindAddr = "localhost:8080"

type config struct{}

func (c *config) Limit() rate.Limit { return 10 }

func (c *config) Burst() int { return 10 }

func (*config) BindAddr() string { return bindAddr }

func testServer(t testing.TB, srv controller.Service) *Controller {
	t.Helper()
	ctrl := &Controller{
		engine:  echo.New(),
		log:     zap.L(),
		cfg:     &config{},
		srv:     srv,
		rateCfg: &config{},
	}
	return ctrl
}
