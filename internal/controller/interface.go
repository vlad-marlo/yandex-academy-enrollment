package controller

import (
	"context"
	mw "github.com/vlad-marlo/yandex-academy-enrollment/internal/middleware"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/datetime"
)

//go:generate mockgen --source=interface.go --destination=mocks/service.go --package=mocks Service

type Config interface {
	BindAddr() string
}

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Service interface {
	GetCourierByID(ctx context.Context, id string) (*model.CourierDTO, error)
	CreateCouriers(ctx context.Context, request *model.CreateCourierRequest) (*model.CouriersCreateResponse, error)
	GetCouriers(ctx context.Context, opts *mw.PaginationOpts) (*model.GetCouriersResponse, error)
	GetCourierMetaInfo(ctx context.Context, req *model.GetCourierMetaInfoRequest) (*model.GetCourierMetaInfoResponse, error)
	GetOrdersAssign(ctx context.Context, date *datetime.Date, id string) (*model.OrderAssignResponse, error)
	GetOrderByID(ctx context.Context, id string) (*model.OrderDTO, error)
	GetOrders(ctx context.Context, opts *mw.PaginationOpts) ([]*model.OrderDTO, error)
	CreateOrders(ctx context.Context, req *model.CreateOrderRequest) ([]*model.OrderDTO, error)
	CompleteOrders(ctx context.Context, req *model.CompleteOrderRequest) ([]*model.OrderDTO, error)
	AssignOrders(ctx context.Context, date *datetime.Date) (*model.OrderAssignResponse, error)
}
