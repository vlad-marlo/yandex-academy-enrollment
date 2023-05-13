package production

import (
	"context"
	"errors"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/store"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"strconv"
)

var _ controller.Service = (*Service)(nil)

func (srv *Service) GetOrdersAssign(ctx context.Context, date *datetime.Date, id string) (order *model.OrderAssignResponse, err error) {
	return nil, ErrNotImplemented
}

func (srv *Service) GetOrderByID(ctx context.Context, id string) (order *model.OrderDTO, err error) {
	var orderID int64
	orderID, err = strconv.ParseInt(id, 10, 0)
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
	orders, err := srv.storage.GetOrders(ctx, opts.Limit(), opts.Offset())
	if err != nil {
		if errors.Is(err, store.ErrNoContent) {
			return []*model.OrderDTO{}, nil
		}
		return nil, ErrNotFound
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
		ord := &model.OrderDTO{
			Weight:        order.Weight,
			Regions:       order.Regions,
			DeliveryHours: order.DeliveryHours,
			Cost:          order.Cost,
		}
		orders = append(orders, ord)
	}
	if err := srv.storage.CreateOrders(ctx, orders); err != nil {
		return nil, ErrBadRequest
	}
	srv.log.Debug("successful created orders")
	return orders, nil
}

func (srv *Service) CompleteOrders(ctx context.Context, req *model.CompleteOrderRequest) ([]*model.OrderDTO, error) {
	if !req.Valid() {
		srv.log.Debug("request didn't pass validation")
		return nil, ErrBadRequest
	}
	return nil, ErrNotImplemented
}

func (srv *Service) AssignOrders(ctx context.Context, date *datetime.Date) (*model.OrderAssignResponse, error) {
	return nil, ErrNotImplemented
}
