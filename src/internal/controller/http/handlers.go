package http

import (
	"github.com/labstack/echo/v4"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"go.uber.org/zap"
	"net/http"
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
	id := c.Param("courier_id")
	courierIDField := zap.String("courier-id", id)

	srv.log.Debug("starting handling getting courier by id", courierIDField)

	courier, err := srv.srv.GetCourierByID(c.Request().Context(), id)
	if err != nil {
		srv.log.Warn("error while getting courier by id", courierIDField, zap.Error(err))
		return srv.checkErr(c, "error while getting courier by id", err)
	}
	srv.log.Debug("successful got courier by id", courierIDField)

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
	opts := GetPaginationOptsFromRequest(c)

	fields := []zap.Field{
		zap.Int("limit", opts.Limit()),
		zap.Int("offset", opts.Offset()),
	}
	srv.log.Debug("handling get couriers", fields...)

	resp, err := srv.srv.GetCouriers(c.Request().Context(), opts)
	if err != nil {
		return srv.checkErr(c, "err while getting response", err, fields...)
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
		return srv.checkErr(c, "err while binding request", err)
	}
	resp, err := srv.srv.CreateCouriers(c.Request().Context(), &request)
	if err != nil {
		return srv.checkErr(c, "err while getting response", err)
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
	var req model.GetCourierMetaInfoRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.BadRequestResponse{})
	}

	resp, err := srv.srv.GetCourierMetaInfo(c.Request().Context(), &req)
	if err != nil {
		return srv.checkErr(c, "err while getting courier meta info", err)
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
	date, err := srv.dateFromContext(c, "date")
	if err != nil {
		return srv.checkErr(c, "err while getting date from context", err)
	}

	var resp *model.OrderAssignResponse

	rawID := c.QueryParam("courier_id")
	resp, err = srv.srv.GetOrdersAssign(c.Request().Context(), date, rawID)
	if err != nil {
		return srv.checkErr(c, "err while getting response", err)
	}

	return c.JSON(http.StatusOK, resp)
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
	id := c.Param("order_id")
	resp, err := srv.srv.GetOrderByID(c.Request().Context(), id)
	if err != nil {
		return srv.checkErr(c, "err while getting courier by id", err)
	}
	return c.JSON(http.StatusOK, resp)
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
	opts := GetPaginationOptsFromRequest(c)
	resp, err := srv.srv.GetOrders(c.Request().Context(), opts)
	if err != nil {
		return srv.checkErr(c, "error while getting orders", err)
	}
	return c.JSON(http.StatusOK, resp)
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
	req := new(model.CreateOrderRequest)
	if err := c.Bind(req); err != nil {
		return srv.checkErr(c, "error while binding request", err)
	}
	resp, err := srv.srv.CreateOrders(c.Request().Context(), req)
	if err != nil {
		return srv.checkErr(c, "error while creating orders", err)
	}
	return c.JSON(http.StatusOK, resp)
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
	req := new(model.CompleteOrderRequest)
	if err := c.Bind(req); err != nil {
		return srv.checkErr(c, "error while binding request", err)
	}
	resp, err := srv.srv.CompleteOrders(c.Request().Context(), req)
	if err != nil {
		return srv.checkErr(c, "error while completing orders", err)
	}
	return c.JSON(http.StatusOK, resp)
}

// HandleAssignOrders assigns orders.
//
//	@Tags		order-controller
//	@Summary	Распределение заказов по курьерам
//	@Accept		json
//	@Produce	json
//	@Param		date	query		string						false	"Дата распределения заказов. Если не указана, то используется текущий день"
//	@Success	201		{object}	model.OrderAssignResponse	"OK"
//	@Failure	400		{object}	model.BadRequestResponse	"Bad Request"
//	@Router		/orders/assign [post]
func (srv *Controller) HandleAssignOrders(c echo.Context) error {
	date, err := srv.dateFromContext(c, "date")
	if err != nil {
		return srv.checkErr(c, "error while getting date from context", err)
	}

	var resp *model.OrderAssignResponse
	resp, err = srv.srv.AssignOrders(c.Request().Context(), date)
	if err != nil {
		return srv.checkErr(c, "error while assigning orders", err)
	}
	return c.JSON(http.StatusCreated, resp)
}
