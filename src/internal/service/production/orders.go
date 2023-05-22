package production

import (
	"context"
	"errors"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/store"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"go.uber.org/zap"
	"strconv"
)

func (srv *Service) AssignOrders(context.Context, *datetime.Date) (*model.OrderAssignResponse, error) {
	return nil, ErrNotImplemented
}

func (srv *Service) GetOrdersAssign(context.Context, *datetime.Date, string) (order *model.OrderAssignResponse, err error) {
	return nil, ErrNotImplemented
}

func (srv *Service) GetOrderByID(ctx context.Context, id string) (order *model.OrderDTO, err error) {
	var orderID int64
	orderID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrBadRequest
	}

	order, err = srv.storage.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, ErrNotFound
	}
	return
}

func (srv *Service) GetOrders(ctx context.Context, opts model.PaginationOpts) ([]*model.OrderDTO, error) {
	if opts == nil {
		return nil, ErrBadRequest
	}

	orders, err := srv.storage.GetOrders(ctx, opts.Limit(), opts.Offset())
	if err != nil {
		if errors.Is(err, store.ErrNoContent) {
			return []*model.OrderDTO{}, nil
		}
		return nil, ErrBadRequest.With(zap.NamedError("storage_error", err))
	}
	return orders, nil
}

func (srv *Service) CreateOrders(ctx context.Context, req *model.CreateOrderRequest) ([]*model.OrderDTO, error) {
	if !req.Valid() {
		srv.log.Debug("request didn't pass validation")
		return nil, ErrBadRequest
	}

	var orders []*model.OrderDTO
	for _, order := range req.Orders {
		orders = append(orders, &model.OrderDTO{
			Weight:        order.Weight,
			Regions:       order.Regions,
			DeliveryHours: order.DeliveryHours,
			Cost:          order.Cost,
		})
	}
	if err := srv.storage.CreateOrders(ctx, orders); err != nil {
		return nil, ErrBadRequest.With(zap.NamedError("storage_error", err))
	}
	srv.log.Debug("successful created orders")
	return orders, nil
}

func (srv *Service) CompleteOrders(ctx context.Context, req *model.CompleteOrderRequest) ([]*model.OrderDTO, error) {
	if !req.Valid() {
		srv.log.Debug("request didn't pass validation")
		return nil, ErrBadRequest
	}

	if err := srv.storage.CompleteOrders(ctx, req.CompleteInfo); err != nil {
		return nil, ErrBadRequest.With(zap.NamedError("storage_error", err))
	}

	var ids []int64
	for _, c := range req.CompleteInfo {
		ids = append(ids, c.OrderID)
	}

	orders, err := srv.storage.GetOrdersByIDs(ctx, ids)
	if err != nil {
		return nil, ErrBadRequest
	}

	return orders, nil
}
