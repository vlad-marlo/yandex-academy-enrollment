package pgx

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/store"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/pgx/client"
	"math/rand"
	"testing"
	"time"
)

func TestStore_GetOrderByID(t *testing.T) {
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, _ := New(cli)

	order, err := s.GetOrderByID(ctx, 1)
	assert.Nil(t, order)
	assert.ErrorIs(t, err, pgx.ErrNoRows)
	want := &model.OrderDTO{
		Weight:  rand.Float64(),
		Regions: rand.Int31(),
		DeliveryHours: []*datetime.TimeInterval{
			datetime.TimeIntervalAlias{Start: 12, End: 24}.TimeInterval(),
		},
		Cost:          rand.Int31(),
		CompletedTime: datetime.Time{},
	}

	require.NoError(t, s.CreateOrders(ctx, []*model.OrderDTO{want}))
	order, err = s.GetOrderByID(ctx, want.OrderID)
	if assert.NotNil(t, order) {
		assert.Equal(t, want, order)
	}
	assert.NoError(t, err)
}

func TestStore_GetOrders(t *testing.T) {
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, err := New(cli)
	require.NoError(t, err)
	orders := []*model.OrderDTO{
		{
			OrderID:       rand.Int63(),
			Weight:        rand.Float64(),
			Regions:       rand.Int31(),
			Cost:          rand.Int31(),
			CompletedTime: datetime.Time{},
			DeliveryHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 123, End: 321}.TimeInterval(),
				datetime.TimeIntervalAlias{Start: 332, End: 400}.TimeInterval(),
			},
		},
		{
			OrderID:       rand.Int63(),
			Weight:        rand.Float64(),
			Regions:       rand.Int31(),
			Cost:          rand.Int31(),
			CompletedTime: datetime.Time{},
			DeliveryHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 11, End: 22}.TimeInterval(),
				datetime.TimeIntervalAlias{Start: 100, End: 200}.TimeInterval(),
			},
		},
		{
			OrderID:       rand.Int63(),
			Weight:        rand.Float64(),
			Regions:       rand.Int31(),
			Cost:          rand.Int31(),
			CompletedTime: datetime.Time{},
			DeliveryHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 0, End: 1200}.TimeInterval(),
			},
		},
	}
	var got []*model.OrderDTO
	err = s.CreateOrders(ctx, orders)
	require.NoError(t, err)

	assert.NotNil(t, orders)
	got, err = s.GetOrders(ctx, 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, orders[0:1], got)
	got, err = s.GetOrders(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, orders[1:2], got)
	got, err = s.GetOrders(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, orders[2:3], got)
	got, err = s.GetOrders(ctx, 2, 2)
	assert.NoError(t, err)
	assert.Equal(t, orders[2:], got)

	got, err = s.GetOrders(ctx, 2, len(orders)+2)
	if err != nil {
		assert.ErrorIs(t, err, store.ErrNoContent)
	}
	assert.NotNil(t, got)
	assert.Empty(t, got)

	got, err = s.GetOrders(ctx, 2, -1)
	assert.Error(t, err)
}

func TestStore_GetCompletedOrdersPriceByCourier_Positive(t *testing.T) {
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, _ := New(cli)

	sum, count, err := s.GetCompletedOrdersPriceByCourier(ctx, 1, time.Unix(0, 0), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, int32(0), sum)
	assert.Equal(t, int32(0), count)
}

func TestStore_GetCompletedOrdersPriceByCourier_Negative_BadCli(t *testing.T) {
	ctx := context.Background()

	cli := client.BadCli(t)

	s, _ := New(cli)

	sum, count, err := s.GetCompletedOrdersPriceByCourier(ctx, 1, time.Unix(0, 0), time.Now())
	assert.Error(t, err)
	assert.Equal(t, int32(0), sum)
	assert.Equal(t, int32(0), count)
}

func TestStore_GetOrdersByIDs_Positive(t *testing.T) {
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, _ := New(cli)

	resp, err := s.GetOrdersByIDs(ctx, []int64{})
	assert.NoError(t, err)
	assert.Empty(t, resp)
	assert.NotNil(t, resp)

	resp, err = s.GetOrdersByIDs(ctx, []int64{1})
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	}

	orders := []*model.OrderDTO{
		{
			OrderID: 0,
			Weight:  rand.Float64(),
			Regions: rand.Int31(),
			DeliveryHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 123, End: 321}.TimeInterval(),
				datetime.TimeIntervalAlias{Start: 332, End: 400}.TimeInterval(),
			},
			Cost: rand.Int31(),
		},
		{
			OrderID:       2,
			Weight:        rand.Float64(),
			Regions:       rand.Int31(),
			DeliveryHours: []*datetime.TimeInterval{},
			Cost:          rand.Int31(),
		},
	}
	require.NoError(t, s.CreateOrders(ctx, orders))

	ids := []int64{1, 2}
	resp, err = s.GetOrdersByIDs(ctx, ids)
	assert.NoError(t, err)
	for i := range orders {
		assert.Equal(t, orders[i], resp[i])
	}
}

func TestGetOrdersByIDs_BadCli(t *testing.T) {
	ctx := context.Background()
	cli := client.BadCli(t)

	s, _ := New(cli)
	resp, err := s.GetOrdersByIDs(ctx, []int64{123})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestStore_CompleteOrders_Positive_NoOrders(t *testing.T) {
	ctx := context.Background()
	cli, td := client.NewTest(t)
	defer td()

	s, _ := New(cli)
	err := s.CompleteOrders(ctx, []model.CompleteOrder{})
	assert.NoError(t, err)
}

func TestStore_CompleteOrders_Negative(t *testing.T) {
	ctx := context.Background()
	cli, td := client.NewTest(t)
	defer td()

	s, _ := New(cli)
	err := s.CompleteOrders(ctx, []model.CompleteOrder{{rand.Int63(), rand.Int63(), datetime.Time{}}})
	assert.NoError(t, err)
}

func TestStore_CompleteOrders_Positive(t *testing.T) {
	ctx := context.Background()
	cli, td := client.NewTest(t)
	defer td()

	s, _ := New(cli)
	err := s.CompleteOrders(ctx, []model.CompleteOrder{{rand.Int63(), rand.Int63(), datetime.Time{}}})
	assert.NoError(t, err)
}
