package controller

import (
	"context"
	"github.com/labstack/echo"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
)

type Interface interface {
}

type Config interface {
	BindAddr() string
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Service interface {
	GetCourierByID(c echo.Context, id int) (*model.Courier, error)
}
