package controller

import (
	"context"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/datetime"
)

//go:generate mockgen --source=service.go --destination=mocks/service.go --package=mocks

type Service interface {
	GetCourierByID(ctx context.Context, id string) (*model.CourierDTO, error)
	CreateCouriers(ctx context.Context, request *model.CreateCourierRequest) (*model.CouriersCreateResponse, error)
	GetCouriers(ctx context.Context, opts model.PaginationOpts) (*model.GetCouriersResponse, error)
	GetCourierMetaInfo(ctx context.Context, req *model.GetCourierMetaInfoRequest) (*model.GetCourierMetaInfoResponse, error)
	GetOrdersAssign(ctx context.Context, date *datetime.Date, id string) (*model.OrderAssignResponse, error)
	GetOrderByID(ctx context.Context, id string) (*model.OrderDTO, error)
	GetOrders(ctx context.Context, opts model.PaginationOpts) ([]*model.OrderDTO, error)
	CreateOrders(ctx context.Context, req *model.CreateOrderRequest) ([]*model.OrderDTO, error)
	CompleteOrders(ctx context.Context, req *model.CompleteOrderRequest) ([]*model.OrderDTO, error)
	AssignOrders(ctx context.Context, date *datetime.Date) (*model.OrderAssignResponse, error)
}
