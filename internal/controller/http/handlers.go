package http

import (
	"github.com/labstack/echo"
	"net/http"
)

func (srv *Controller) HandlePing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (srv *Controller) HandleGetCourier(c echo.Context) error {
	return nil
}
