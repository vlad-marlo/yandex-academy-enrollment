package http

import (
	"github.com/labstack/echo/v4"
	mw "github.com/vlad-marlo/yandex-academy-enrollment/internal/middleware"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"net/http"
	"strconv"
)

func (srv *Controller) HandlePing(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

// HandleGetCourier return courier with provided id.
//
//	@Tags		courier-controller
//	@Summary	Получение профиля курьера
//	@Accept		json
//	@Produce	json
//	@Param		courier_id	path		int							true	"Courier identifier"
//	@Success	200			{object}	model.CourierDTO			"OK"
//	@Failure	400			{object}	model.BadRequestResponse	"Bad Request"
//	@Failure	404			{object}	model.BadRequestResponse	"Not Found"
//	@Router		/couriers/{courier_id} [get]
func (srv *Controller) HandleGetCourier(c echo.Context) error {
	var courier *model.CourierDTO

	sID := c.Param("courier_id")
	id, err := strconv.Atoi(sID)
	if err != nil {
		return c.JSON(http.StatusNotFound, nil)
	}
	courier, err = srv.srv.GetCourierByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, courier)
}

// HandleGetCouriers return courier with provided id.
//
//	@Tags		courier-controller
//	@Summary	Получение профилей курьеров
//	@Accept		json
//	@Produce	json
//	@Param		limit	query		int							false	"Максимальное количество курьеров в выдаче. Если параметр не передан, то значение по умолчанию равно 1."
//	@Param		offset	query		int							false	"Количество курьеров, которое нужно пропустить для отображения текущей страницы. Если параметр не передан, то значение по умолчанию равно 0."
//	@Success	200		{object}	model.GetCouriersResponse	"OK"
//	@Failure	400		{object}	model.BadRequestResponse	"Bad Request"
//	@Router		/couriers/ [get]
func (srv *Controller) HandleGetCouriers(c echo.Context) error {
	opts := mw.GetPaginationOptsFromRequest(c)
	resp, err := srv.srv.GetCouriers(c.Request().Context(), opts)
	if err != nil {
		return srv.checkErr(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

// HandleCreateCouriers creates provided couriers.
//
//	@Tags		courier-controller
//	@Summary	Создание профилей курьеров
//	@Accept		json
//	@Produce	json
//	@Param		request	body		model.CreateCourierRequest		true	"Couriers"
//	@Success	200		{object}	model.CouriersCreateResponse	"OK"
//	@Failure	400		{object}	model.BadRequestResponse		"Bad Request"
//	@Router		/couriers/ [post]
func (srv *Controller) HandleCreateCouriers(c echo.Context) error {
	var request model.CreateCourierRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp, err := srv.srv.CreateCouriers(c.Request().Context(), &request)
	if err != nil {
		return srv.checkErr(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

// HandleGetCourierMetaInfo return courier meta info.
//
//	@Tags		courier-controller
//	@Summary	Получение meta-информации о курьере.
//	@Accept		json
//	@Produce	json
//	@Param		courier_id	path		int									true	"Courier identifier"
//	@Param		startDate	query		string								true	"Максимальное количество курьеров в выдаче. Если параметр не передан, то значение по умолчанию равно 1."
//	@Param		endDate		query		string								true	"Количество курьеров, которое нужно пропустить для отображения текущей страницы. Если параметр не передан, то значение по умолчанию равно 0."
//	@Success	200			{object}	model.GetCourierMetaInfoResponse	"OK"
//	@Failure	400			{object}	model.BadRequestResponse			"Bad Request"
//	@Router		/couriers/meta-info/{courier_id} [get]
func (srv *Controller) HandleGetCourierMetaInfo(c echo.Context) error {
	var req *model.GetCourierMetaInfoRequest
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	resp, err := srv.srv.GetCourierMetaInfo(c.Request().Context(), req)
	if err != nil {
		return srv.checkErr(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

// HandleGetOrdersAssign doc.
//
//	@Tags		courier-controller
//	@Summary	список распределенных заказов
//	@Accept		json
//	@Produce	json
//	@Param		courier_id	query		int							false	"Идентификатор курьера для получения списка распредленных заказов. Если не указан, возвращаются данные по всем курьерам."
//	@Param		date		query		string						false	"Дата распределения заказов. Если не указана, то используется текущий день"
//	@Success	200			{object}	model.OrderAssignResponse	"OK"
//	@Failure	400			{object}	model.BadRequestResponse	"Bad Request"
//	@Router		/couriers/assignments [get]
func (srv *Controller) HandleGetOrdersAssign(c echo.Context) error {
	return nil
}

// HandleGetOrder return courier with provided id.
//
//	@Tags		order-controller
//	@Summary	Получение информации о заказе
//	@Accept		json
//	@Produce	json
//	@Param		order_id	path		int							true	"Order identifier"
//	@Success	200			{object}	model.OrderDTO				"OK"
//	@Failure	400			{object}	model.BadRequestResponse	"Bad Request"
//	@Failure	404			{object}	model.BadRequestResponse	"Not Found"
//	@Router		/orders/{order_id} [get]
func (srv *Controller) HandleGetOrder(c echo.Context) error {
	panic("not implemented")
}

// HandleGetOrders return courier with provided id.
//
//	@Tags		order-controller
//	@Summary	Получение заказов
//	@Accept		json
//	@Produce	json
//	@Param		limit	query		int							false	"Максимальное количество заказов в выдаче. Если параметр не передан, то значение по умолчанию равно 1."
//	@Param		offset	query		int							false	"Количество заказов, которое нужно пропустить для отображения текущей страницы. Если параметр не передан, то значение по умолчанию равно 0."
//	@Success	200		{array}		model.OrderDTO				"OK"
//	@Failure	400		{object}	model.BadRequestResponse	"Bad Request"
//	@Router		/orders/ [get]
func (srv *Controller) HandleGetOrders(c echo.Context) error {
	panic("not implemented")
}

// HandleCreateOrders creates provided orders.
//
//	@Tags		order-controller
//	@Summary	Создание заказов
//	@Accept		json
//	@Produce	json
//	@Param		request	body		model.CreateOrderRequest	true	"Orders"
//	@Success	200		{array}		model.OrderDTO				"OK"
//	@Failure	400		{object}	model.BadRequestResponse	"Bad Request"
//	@Router		/orders/ [post]
func (srv *Controller) HandleCreateOrders(c echo.Context) error {
	panic("not implemented")
}

// HandleCompleteOrders completes provided orders.
//
// This handler is idempotent.
//
//	@Tags		order-controller
//	@Summary	Завершение заказов
//	@Accept		json
//	@Produce	json
//	@Param		request	body		model.CompleteOrderRequest	true	"Orders"
//	@Success	200		{array}		model.OrderDTO				"OK"
//	@Failure	400		{object}	model.BadRequestResponse	"Bad Request"
//	@Router		/orders/complete [post]
func (srv *Controller) HandleCompleteOrders(c echo.Context) error {
	panic("not implemented")
}

// HandleAssignOrders assigns orders.
//
//	@Tags		order-controller
//	@Summary	Распределение заказов по курьерам
//	@Accept		json
//	@Produce	json
//	@Param		date	query		string						false	"Дата распределения заказов. Если не указана, то используется текущий день"
//	@Success	200		{object}	model.OrderAssignResponse	"OK"
//	@Failure	400		{object}	model.BadRequestResponse	"Bad Request"
//	@Router		/orders/assign [post]
func (srv *Controller) HandleAssignOrders(c echo.Context) error {
	panic("not implemented")
}
