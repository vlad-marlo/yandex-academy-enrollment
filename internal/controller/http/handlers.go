package http

import (
	"github.com/labstack/echo/v4"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"net/http"
	"strconv"
)

func (srv *Controller) HandlePing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

// HandleGetCourier returns courier with provided id.
//
//	@Tags		User
//	@Summary	Создание пользователя
//	@ID			courier_get
//	@Accept		json
//	@Produce	json
//	@Param		courier_id	path		int	true	"CourierDTO identifier"
//	@Success	201			{object}	model.CouriersCreateResponse
//	@Failure	400			{object}	model.Error
//	@Router		/couriers/{courier_id} [get]
func (srv *Controller) HandleGetCourier(c echo.Context) error {
	var courier *model.CourierDTO

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
