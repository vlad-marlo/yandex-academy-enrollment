package http

import (
	"github.com/labstack/echo"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"net/http"
	"strconv"
)

func (srv *Controller) HandlePing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (srv *Controller) HandleGetCourier(c echo.Context) error {
	var courier *model.Courier

	sID := c.Param("courier_id")
	id, err := strconv.Atoi(sID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	courier, err = srv.srv.GetCourierByID(c, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, courier)
}
