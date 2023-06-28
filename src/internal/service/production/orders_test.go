package production

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller/http"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/service/production/mocks"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/store"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"math/rand"
	"testing"
	"time"
)

func TestService_GetOrdersAssign(t *testing.T) {
	srv := testService(t, nil)
	resp, err := srv.GetOrdersAssign(context.Background(), nil, "")
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNotImplemented)
	}
}

func TestService_AssignOrders(t *testing.T) {
	srv := testService(t, nil)
	resp, err := srv.AssignOrders(context.Background(), nil)
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNotImplemented)
	}
}

func TestService_GetOrderByID_Negative_UnParsableID(t *testing.T) {
	tt := []string{"random string", "1.1", "1,1", ""}
	for _, tc := range tt {
		t.Run(tc, func(t *testing.T) {
			srv := testService(t, nil)
			resp, err := srv.GetOrderByID(context.Background(), tc)
			assert.Nil(t, resp)
			if assert.Error(t, err) {
				assert.ErrorIs(t, err, ErrBadRequest)
			}
		})
	}
}

func TestService_GetOrderByID_Negative_ErrInStorage(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	str := mocks.NewMockStore(ctrl)
	str.EXPECT().GetOrderByID(ctx, int64(123)).Return(nil, store.ErrNoContent)
	srv := testService(t, str)

	resp, err := srv.GetOrderByID(ctx, "123")
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNotFound)
	}
}

func TestService_GetOrderByID_Positive_NilResp(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	str := mocks.NewMockStore(ctrl)
	str.EXPECT().GetOrderByID(ctx, int64(123)).Return(nil, nil)
	srv := testService(t, str)

	resp, err := srv.GetOrderByID(ctx, "123")
	assert.Nil(t, resp)
	assert.NoError(t, err)
}

func TestService_GetOrderByID_Positive_NonNilResp(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	str := mocks.NewMockStore(ctrl)

	want := &model.OrderDTO{
		OrderID:       rand.Int63(),
		Weight:        rand.Float64(),
		Regions:       rand.Int31(),
		DeliveryHours: []*datetime.TimeInterval{},
		Cost:          rand.Int31(),
		CompletedTime: datetime.Time{},
	}
	str.EXPECT().GetOrderByID(ctx, int64(123)).Return(want, nil)

	srv := testService(t, str)

	resp, err := srv.GetOrderByID(ctx, "123")
	if assert.NotNil(t, resp) {
		assert.Equal(t, want, resp)
	}
	assert.NoError(t, err)
}

func TestService_GetOrders_Negative_NilReference(t *testing.T) {
	srv := testService(t, nil)
	resp, err := srv.GetOrders(context.Background(), nil)
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrBadRequest)
	}
}
func TestService_GetOrders_Negative_ErrInStorage(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)

	str.EXPECT().GetOrders(ctx, 1, 0).Return(nil, errors.New(""))

	srv := testService(t, str)
	resp, err := srv.GetOrders(ctx, http.NewPaginationOpts("", ""))
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrBadRequest)
	}
}

func TestService_GetOrders_Positive_NoContentError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)

	str.EXPECT().GetOrders(ctx, 1, 0).Return(nil, store.ErrNoContent)

	srv := testService(t, str)
	resp, err := srv.GetOrders(ctx, http.NewPaginationOpts("", ""))
	if assert.NotNil(t, resp) {
		assert.Empty(t, resp)
	}
	assert.NoError(t, err)
}

func TestService_GetOrders_Positive_1(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)

	str.EXPECT().GetOrders(ctx, 1, 0).Return(nil, nil)

	srv := testService(t, str)
	resp, err := srv.GetOrders(ctx, http.NewPaginationOpts("", ""))
	assert.Nil(t, resp)
	assert.NoError(t, err)
}

func TestService_GetOrders_Positive_2(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)

	str.EXPECT().GetOrders(ctx, 1, 0).Return([]*model.OrderDTO{
		{
			OrderID:       1,
			Weight:        2,
			Regions:       3,
			DeliveryHours: nil,
			Cost:          3,
			CompletedTime: datetime.Time{},
		},
		new(model.OrderDTO),
	}, nil)

	srv := testService(t, str)
	resp, err := srv.GetOrders(ctx, http.NewPaginationOpts("", ""))
	if assert.NotNil(t, resp) {
		assert.NotEmpty(t, resp)
	}
	assert.NoError(t, err)
}

func TestService_CreateOrders_Negative_NonValid(t *testing.T) {
	srv := testService(t, nil)
	resp, err := srv.CreateOrders(context.Background(), nil)
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrBadRequest)
	}
}

func TestService_CreateOrders_Negative_StoreError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	req := &model.CreateOrderRequest{
		Orders: []model.CreateOrderDTO{
			{
				Weight:  rand.Float64(),
				Regions: rand.Int31(),
				DeliveryHours: []*datetime.TimeInterval{
					datetime.TimeIntervalAlias{Start: 12, End: 33}.TimeInterval(),
					datetime.TimeIntervalAlias{Start: 11, End: 34}.TimeInterval(),
				},
				Cost: rand.Int31(),
			},
		},
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

	str.EXPECT().CreateOrders(ctx, orders).Return(store.ErrNoContent)

	resp, err := srv.CreateOrders(ctx, req)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrBadRequest)
	}
	assert.Nil(t, resp)
}

func TestService_CreateOrders_Positive(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	req := &model.CreateOrderRequest{
		Orders: []model.CreateOrderDTO{
			{
				Weight:  rand.Float64(),
				Regions: rand.Int31(),
				DeliveryHours: []*datetime.TimeInterval{
					datetime.TimeIntervalAlias{Start: 12, End: 33}.TimeInterval(),
					datetime.TimeIntervalAlias{Start: 11, End: 34}.TimeInterval(),
				},
				Cost: rand.Int31(),
			},
		},
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

	str.EXPECT().CreateOrders(ctx, orders).Return(nil)

	resp, err := srv.CreateOrders(ctx, req)
	if assert.NotNil(t, resp) {
		assert.Equal(t, orders, resp)
	}
	assert.NoError(t, err)
}

func TestService_CompleteOrders_NegativeBadRequest(t *testing.T) {
	srv := testService(t, nil)
	resp, err := srv.CompleteOrders(context.Background(), nil)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrBadRequest)
	}
	assert.Nil(t, resp)
}

func TestService_CompleteOrders_Negative_CompleteOrders(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	req := &model.CompleteOrderRequest{
		CompleteInfo: []model.CompleteOrder{
			{
				CourierID:    123,
				OrderID:      321,
				CompleteTime: datetime.Time(time.Now()),
			},
		},
	}

	str.EXPECT().CompleteOrders(ctx, req.CompleteInfo).Return(errors.New(""))

	resp, err := srv.CompleteOrders(ctx, req)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrBadRequest)
	}
	assert.Nil(t, resp)
}

func TestService_CompleteOrders_Negative_GetOrders(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	req := &model.CompleteOrderRequest{
		CompleteInfo: []model.CompleteOrder{
			{
				CourierID:    123,
				OrderID:      321,
				CompleteTime: datetime.Time(time.Now()),
			},
		},
	}

	str.EXPECT().CompleteOrders(ctx, req.CompleteInfo).Return(nil)
	str.EXPECT().GetOrdersByIDs(ctx, []int64{321}).Return(nil, errors.New(""))

	resp, err := srv.CompleteOrders(ctx, req)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrBadRequest)
	}
	assert.Nil(t, resp)
}

func TestService_CompleteOrders_Positive(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	req := &model.CompleteOrderRequest{
		CompleteInfo: []model.CompleteOrder{
			{
				CourierID:    123,
				OrderID:      321,
				CompleteTime: datetime.Time(time.Now()),
			},
		},
	}
	expected := []*model.OrderDTO{
		{
			OrderID:       req.CompleteInfo[0].OrderID,
			Weight:        rand.Float64(),
			Regions:       rand.Int31(),
			DeliveryHours: []*datetime.TimeInterval{},
			Cost:          rand.Int31(),
			CompletedTime: req.CompleteInfo[0].CompleteTime,
		},
	}

	str.EXPECT().CompleteOrders(ctx, req.CompleteInfo).Return(nil)
	str.EXPECT().GetOrdersByIDs(ctx, []int64{321}).Return(expected, nil)

	resp, err := srv.CompleteOrders(ctx, req)
	assert.NoError(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, expected, resp)
	}
}
