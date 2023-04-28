package http

import (
	"github.com/labstack/echo"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (srv *Controller) HandlePing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (srv *Controller) HandleGetCourier(c echo.Context) error {
	var courier *model.Courier

	defer func() {
		if err := c.Request().Body.Close(); err != nil {
			srv.log.Warn(
				"error was occurred while closing request body in get courier by id handler",
				zap.Error(err),
			)
		}
	}()

	sID := c.Param("courier_id")
	id, err := strconv.Atoi(sID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	courier, err = srv.srv.GetCourierByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, courier)
}

func (srv *Controller) HandleCreateCouriers(c echo.Context) error {
	var request model.CouriersCreateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp, err := srv.srv.CreateCouriers(c.Request().Context(), &request)
	if err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, resp)
}
