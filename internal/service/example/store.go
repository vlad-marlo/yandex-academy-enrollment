package example

import (
	"context"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/fielderr"
	"math/rand"
	"time"
)

var (
	ErrBadRequest = fielderr.New("bad request", model.BadRequestResponse{}, fielderr.CodeBadRequest)
	timeInterval1 = datetime.TimeIntervalAlias{Start: 4 * 60, End: 20 * 60}.TimeInterval()
	timeInterval2 = datetime.TimeIntervalAlias{Start: 20 * 60, End: 4 * 60}.TimeInterval()
	timeInterval3 = datetime.TimeIntervalAlias{Start: 12 * 60, End: 15 * 60}.TimeInterval()
	couriers      = []model.CourierDTO{
		{
			CourierID:   1,
			CourierType: model.AutoCourierTypeString,
			Regions:     []int{799, 77},
			WorkingHours: []*datetime.TimeInterval{
				timeInterval1,
			},
		},
		{
			CourierID:   2,
			CourierType: model.BikeCourierTypeString,
			Regions:     []int{4, 91},
			WorkingHours: []*datetime.TimeInterval{
				timeInterval2,
			},
		},
		{
			CourierID:   3,
			CourierType: model.FootCourierTypeString,
			Regions:     []int{29},
			WorkingHours: []*datetime.TimeInterval{
				timeInterval3,
			},
		},
	}
	couriersMetaInfo = []*model.GetCourierMetaInfoResponse{
		{
			CourierID:   1,
			CourierType: model.AutoCourierTypeString,
			Regions:     []int{799, 77},
			WorkingHours: []*datetime.TimeInterval{
				timeInterval1,
			},
			Rating:   23,
			Earnings: 0,
		},
		{
			CourierID:   2,
			CourierType: model.BikeCourierTypeString,
			Regions:     []int{4, 91},
			WorkingHours: []*datetime.TimeInterval{
				timeInterval2,
			},
			Rating:   3321312,
			Earnings: 9133,
		},
		{
			CourierID:   3,
			CourierType: model.FootCourierTypeString,
			Regions:     []int{29},
			WorkingHours: []*datetime.TimeInterval{
				timeInterval3,
			},
			Rating:   13,
			Earnings: 7873,
		},
	}
	orders = []model.OrderDTO{
		{
			OrderID: 1,
			Weight:  rand.Float64(),
			Regions: 77,
			DeliveryHours: []*datetime.TimeInterval{
				timeInterval1,
				timeInterval2,
			},
			Cost:          rand.Int(),
			CompletedTime: datetime.TimeAlias(time.Now()),
		},
		{
			OrderID: 2,
			Weight:  rand.Float64(),
			Regions: 4,
			DeliveryHours: []*datetime.TimeInterval{
				timeInterval1,
				timeInterval2,
			},
			Cost:          rand.Int(),
			CompletedTime: datetime.TimeAlias(time.Now()),
		},
		{
			OrderID: 3,
			Weight:  rand.Float64(),
			Regions: 29,
			DeliveryHours: []*datetime.TimeInterval{
				timeInterval1,
				timeInterval3,
			},
			Cost:          rand.Int(),
			CompletedTime: datetime.TimeAlias(time.Now()),
		},
	}
)

var _ controller.Service = (*service)(nil)

type service struct{}

func New() controller.Service {
	return service{}
}

func (s service) GetCourierByID(context.Context, string) (*model.CourierDTO, error) {
	index := rand.Int() % len(couriers)
	return &couriers[index], nil
}

func (s service) CreateCouriers(context.Context, *model.CreateCourierRequest) (*model.CouriersCreateResponse, error) {
	return &model.CouriersCreateResponse{
		Couriers: couriers,
	}, nil
}

func (s service) GetCouriers(_ context.Context, opts model.PaginationOpts) (res *model.GetCouriersResponse, err error) {
	res = new(model.GetCouriersResponse)
	res.Couriers = []model.CourierDTO{}
	res.Limit = opts.Limit()
	res.Offset = opts.Offset()
	if opts.Offset() > len(couriers) {
		return
	}
	if opts.Offset()+opts.Limit() > len(couriers) {
		for i := opts.Offset(); i < len(couriers); i++ {
			res.Couriers = append(res.Couriers, couriers[i])
		}
		return
	}
	for i := opts.Offset(); i < opts.Offset()+opts.Limit(); i++ {
		res.Couriers = append(res.Couriers, couriers[i])
	}
	return
}

func (s service) GetCourierMetaInfo(context.Context, *model.GetCourierMetaInfoRequest) (*model.GetCourierMetaInfoResponse, error) {
	idx := rand.Int() % len(couriersMetaInfo)
	return couriersMetaInfo[idx], nil
}

func (s service) GetOrdersAssign(_ context.Context, date *datetime.Date, _ string) (*model.OrderAssignResponse, error) {
	return &model.OrderAssignResponse{
		Date: date.String(),
		Couriers: []model.CourierGroupOrders{
			{
				CourierID: 123,
				Orders: []model.GroupOrders{
					{
						GroupOrderID: 1,
						Orders:       orders,
					},
				},
			},
			{
				CourierID: 321,
				Orders:    []model.GroupOrders{},
			},
		},
	}, nil
}

func (s service) GetOrderByID(context.Context, string) (*model.OrderDTO, error) {
	return &orders[rand.Int()%len(orders)], nil
}

func (s service) GetOrders(context.Context, model.PaginationOpts) (res []*model.OrderDTO, err error) {
	for _, i := range orders {
		res = append(res, &model.OrderDTO{
			OrderID:       i.OrderID,
			Weight:        i.Weight,
			Regions:       i.Regions,
			DeliveryHours: i.DeliveryHours,
			Cost:          i.Cost,
			CompletedTime: i.CompletedTime,
		})
	}
	return
}

func (s service) CreateOrders(_ context.Context, req *model.CreateOrderRequest) (res []*model.OrderDTO, err error) {
	if !req.Valid() {
		return nil, ErrBadRequest
	}
	for _, o := range req.Orders {
		res = append(res, &model.OrderDTO{
			OrderID:       rand.Int(),
			Weight:        o.Weight,
			Regions:       o.Regions,
			DeliveryHours: o.DeliveryHours,
			Cost:          o.Cost,
			CompletedTime: datetime.TimeAlias(time.Now()),
		})
	}
	return
}

func (s service) CompleteOrders(_ context.Context, req *model.CompleteOrderRequest) (res []*model.OrderDTO, err error) {
	if !req.Valid() {
		return nil, ErrBadRequest
	}
	for _, completeOrder := range req.CompleteInfo {
		res = append(res, &model.OrderDTO{
			OrderID: completeOrder.OrderID,
			Weight:  rand.Float64(),
			Regions: rand.Int(),
			DeliveryHours: []*datetime.TimeInterval{
				timeInterval3,
			},
			Cost:          rand.Int(),
			CompletedTime: completeOrder.CompleteTime,
		})
	}
	return
}

func (s service) AssignOrders(_ context.Context, date *datetime.Date) (*model.OrderAssignResponse, error) {
	return &model.OrderAssignResponse{
		Date: date.String(),
		Couriers: []model.CourierGroupOrders{
			{
				CourierID: 123,
				Orders: []model.GroupOrders{
					{
						GroupOrderID: 1,
						Orders:       orders,
					},
				},
			},
			{
				CourierID: 321,
				Orders:    []model.GroupOrders{},
			},
		},
	}, nil
}
