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
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/pgx/migrator"
	"testing"
)

func TestStore_getCourierRegions_PositiveNoData(t *testing.T) {
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, err := New(cli)
	require.NoError(t, err)

	var regions []int32
	regions, err = s.getCourierRegions(ctx, 0)

	assert.NoError(t, err)
	if assert.NotNil(t, regions) {
		assert.Empty(t, regions)
	}
}

func TestStore_getCourierRegions_Negative(t *testing.T) {
	ctx := context.Background()

	cli := client.BadCli(t)

	s, err := New(cli)
	require.NoError(t, err)

	var regions []int32
	regions, err = s.getCourierRegions(ctx, 0)

	assert.Error(t, err)
	t.Log(err)
	if assert.Nil(t, regions) {
		assert.Empty(t, regions)
	}
}

func TestStore_getCourierWorkingHours_PositiveNoData(t *testing.T) {
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, err := New(cli)
	require.NoError(t, err)

	var hours []*datetime.TimeInterval
	hours, err = s.getCourierWorkingHours(ctx, 0)

	assert.NoError(t, err)
	if assert.NotNil(t, hours) {
		assert.Empty(t, hours)
	}
}

func TestStore_getCourierWorkingHours_Negative(t *testing.T) {
	ctx := context.Background()

	cli := client.BadCli(t)

	s, err := New(cli)
	require.NoError(t, err)

	var hours []*datetime.TimeInterval
	hours, err = s.getCourierWorkingHours(ctx, 0)

	assert.Error(t, err)
	t.Log(err)
	if assert.Nil(t, hours) {
		assert.Empty(t, hours)
	}
}

func TestStore_GetCourierByID_NegativeNotFound(t *testing.T) {
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, err := New(cli)
	require.NoError(t, err)

	var courier *model.CourierDTO
	courier, err = s.GetCourierByID(ctx, 0)
	assert.Nil(t, courier)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	}
}

func TestStore_addRegionsToCourier_Negative_NonExists(t *testing.T) {
	var (
		tx pgx.Tx
	)
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, err := New(cli)
	require.NoError(t, err)

	tx, err = cli.P().Begin(ctx)
	defer assert.NoError(t, tx.Rollback(ctx))
	require.NoError(t, err)

	err = s.addRegionsToCourier(ctx, tx, model.CourierDTO{CourierID: 1, Regions: []int32{1, 2, 3}})
	assert.Error(t, err)
}

func TestStore_addWorkingHoursToCourier_Negative_NonExists(t *testing.T) {
	var (
		tx pgx.Tx
	)
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, err := New(cli)
	require.NoError(t, err)

	tx, err = cli.P().Begin(ctx)
	defer assert.NoError(t, tx.Rollback(ctx))
	require.NoError(t, err)

	err = s.addWorkingHoursToCourier(ctx, tx, model.CourierDTO{CourierID: 1, WorkingHours: []*datetime.TimeInterval{
		datetime.TimeIntervalAlias{Start: 123, End: 321}.TimeInterval(),
	}})
	assert.Error(t, err)
}

func TestStore_CreateCouriers_Positive(t *testing.T) {
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, err := New(cli)
	require.NoError(t, err)
	var c *model.CourierDTO
	c, err = s.GetCourierByID(ctx, 1)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	}
	assert.Nil(t, c)

	couriers := []model.CreateCourierDTO{
		{
			CourierType: model.BikeCourierTypeString,
			Regions:     []int32{1, 3},
			WorkingHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 123, End: 321}.TimeInterval(),
				datetime.TimeIntervalAlias{Start: 332, End: 400}.TimeInterval(),
			},
		},
	}

	var resp []model.CourierDTO

	resp, err = s.CreateCouriers(ctx, couriers)
	assert.NoError(t, err)
	if assert.NotNilf(t, resp, "err = %v, resp=%v", err, resp) {
		require.NotEmpty(t, resp)
		courier := resp[0]
		var (
			wh      []*datetime.TimeInterval
			regions []int32
			c       *model.CourierDTO
		)
		assert.Equal(t, int64(1), courier.CourierID)
		wh, err = s.getCourierWorkingHours(ctx, courier.CourierID)
		assert.NoErrorf(t, err, "unxepectedly got non nil error: %v", err)
		assert.Equal(t, courier.WorkingHours, wh)

		regions, err = s.getCourierRegions(ctx, courier.CourierID)
		assert.NoError(t, err)
		assert.Equal(t, courier.Regions, regions)

		c, err = s.GetCourierByID(ctx, courier.CourierID)
		assert.Equal(t, courier, *c)
	}
}

func TestStore_CreateCouriers_Negative_BadCli(t *testing.T) {
	ctx := context.Background()

	cli := client.BadCli(t)

	s, err := New(cli)
	require.NoError(t, err)
	couriers := []model.CreateCourierDTO{
		{
			CourierType: model.BikeCourierTypeString,
			Regions:     []int32{1, 3},
			WorkingHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 123, End: 321}.TimeInterval(),
				datetime.TimeIntervalAlias{Start: 332, End: 400}.TimeInterval(),
			},
		},
	}
	var resp []model.CourierDTO
	resp, err = s.CreateCouriers(ctx, couriers)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestStore_CreateCouriers_Negative_BadData(t *testing.T) {
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, err := New(cli)
	require.NoError(t, err)
	couriers := []model.CreateCourierDTO{
		{
			CourierType: "unknown type",
			Regions:     []int32{1, 3},
			WorkingHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 123, End: 321}.TimeInterval(),
				datetime.TimeIntervalAlias{Start: 332, End: 400}.TimeInterval(),
			},
		},
	}
	var resp []model.CourierDTO
	resp, err = s.CreateCouriers(ctx, couriers)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestStore_createCourier_Negative(t *testing.T) {
	if testing.Short() {
		return
	}
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, err := New(cli)
	require.NoError(t, err)

	var tx pgx.Tx
	tx, err = cli.P().Begin(ctx)
	assert.NoError(t, err)

	couriers := model.CreateCourierDTO{
		CourierType: model.BikeCourierTypeString,
		Regions:     []int32{1, 3},
		WorkingHours: []*datetime.TimeInterval{
			datetime.TimeIntervalAlias{Start: 123, End: 321}.TimeInterval(),
			datetime.TimeIntervalAlias{Start: 332, End: 400}.TimeInterval(),
		},
	}
	i, err := migrator.MigrateDown(cli)
	t.Log(i)
	require.NoError(t, err)
	var resp model.CourierDTO
	resp, err = s.createCourier(ctx, tx, couriers)
	assert.Error(t, err)
	assert.Empty(t, resp)
}

func TestStore_GetCouriers(t *testing.T) {
	ctx := context.Background()

	cli, td := client.NewTest(t)
	defer td()

	s, err := New(cli)
	require.NoError(t, err)
	couriers := []model.CreateCourierDTO{
		{
			CourierType: model.BikeCourierTypeString,
			Regions:     []int32{1, 3},
			WorkingHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 123, End: 321}.TimeInterval(),
				datetime.TimeIntervalAlias{Start: 332, End: 400}.TimeInterval(),
			},
		},
		{
			CourierType: model.AutoCourierTypeString,
			Regions:     []int32{4, 12, 33},
			WorkingHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 11, End: 22}.TimeInterval(),
				datetime.TimeIntervalAlias{Start: 100, End: 200}.TimeInterval(),
			},
		},
		{
			CourierType: model.FootCourierTypeString,
			Regions:     []int32{6},
			WorkingHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 0, End: 1200}.TimeInterval(),
			},
		},
	}
	var resp, got []model.CourierDTO
	resp, err = s.CreateCouriers(ctx, couriers)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	got, err = s.GetCouriers(ctx, 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, resp[0:1], got)
	got, err = s.GetCouriers(ctx, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, resp[1:2], got)
	got, err = s.GetCouriers(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, resp[2:3], got)
	got, err = s.GetCouriers(ctx, 2, 2)
	assert.NoError(t, err)
	assert.Equal(t, resp[2:], got)

	got, err = s.GetCouriers(ctx, 2, len(resp)+2)
	if err != nil {
		assert.ErrorIs(t, err, store.ErrNoContent)
	}
	assert.NotNil(t, got)
	assert.Empty(t, got)

	got, err = s.GetCouriers(ctx, 2, -1)
	assert.Error(t, err)
}

func TestStore_GetCouriers_NegativeBadCli(t *testing.T) {
	cli := client.BadCli(t)
	s, err := New(cli)
	require.NoError(t, err)
	resp, err := s.GetCouriers(context.Background(), 0, 0)
	assert.Nil(t, resp)
	assert.Error(t, err)
}
