package model

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/datetime"
	"testing"
)

const (
	testTimeInterval1String = "12:00-13:30"
	testTimeInterval2String = "13:50-16:50"
)

var (
	testTimeInterval1 = func(t testing.TB) *datetime.TimeInterval {
		t.Helper()
		ti, err := datetime.ParseTimeInterval(testTimeInterval1String)
		require.NoError(t, err)
		return ti
	}
	testTimeInterval2 = func(t testing.TB) *datetime.TimeInterval {
		t.Helper()
		ti, err := datetime.ParseTimeInterval(testTimeInterval2String)
		require.NoError(t, err)
		return ti
	}
)

func TestAll(t *testing.T) {
	tt := []struct {
		name  string
		items []any
		f     func(any) bool
		want  assert.BoolAssertionFunc
	}{
		{"positive #1", []any{1, ""}, func(any) bool { return true }, assert.True},
		{"negative #1", []any{1, ""}, func(any) bool { return false }, assert.False},
		{
			"negative #2",
			[]any{1, ""},
			func(i any) bool {
				_, ok := i.(string)
				return ok
			},
			assert.False,
		},
		{
			"positive #2",
			[]any{"string 1", "string 2"},
			func(i any) bool {
				_, ok := i.(string)
				return ok
			},
			assert.True,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.want(t, all(tc.f, tc.items))
		})
	}
}

func TestCompleteOrderRequest_Valid_Negative_NilReference(t *testing.T) {
	req := (*CompleteOrderRequest)(nil)
	assert.False(t, req.Valid())
}

func TestCompleteOrderRequest_Valid(t *testing.T) {
	tt := []struct {
		name string
		req  *CompleteOrderRequest
		want assert.BoolAssertionFunc
	}{
		{"negative #1 - nil reference", (*CompleteOrderRequest)(nil), assert.False},
		{"negative #2 - no orders", new(CompleteOrderRequest), assert.False},
		{
			"positive #1",
			&CompleteOrderRequest{
				CompleteInfo: []CompleteOrder{
					{
						CourierID:    1,
						OrderID:      2,
						CompleteTime: datetime.TimeAlias{},
					},
					{
						CourierID:    2,
						OrderID:      1,
						CompleteTime: datetime.TimeAlias{},
					},
				},
			},
			assert.True,
		},
		{
			"positive #1",
			&CompleteOrderRequest{
				CompleteInfo: []CompleteOrder{
					{
						CourierID:    1,
						OrderID:      2,
						CompleteTime: datetime.TimeAlias{},
					},
					{
						CourierID:    2,
						OrderID:      1,
						CompleteTime: datetime.TimeAlias{},
					},
				},
			},
			assert.True,
		},
		{
			"positive #2",
			&CompleteOrderRequest{
				CompleteInfo: []CompleteOrder{
					{
						CourierID:    1,
						OrderID:      2,
						CompleteTime: datetime.TimeAlias{},
					},
				},
			},
			assert.True,
		},
		{
			"negative #1",
			&CompleteOrderRequest{
				CompleteInfo: []CompleteOrder{
					{
						CourierID:    1,
						OrderID:      1,
						CompleteTime: datetime.TimeAlias{},
					},
					{
						CourierID:    2,
						OrderID:      1,
						CompleteTime: datetime.TimeAlias{},
					},
				},
			},

			assert.False,
		},
		{
			"negative #2",
			&CompleteOrderRequest{
				CompleteInfo: []CompleteOrder{
					{
						CourierID:    1,
						OrderID:      2,
						CompleteTime: datetime.TimeAlias{},
					},
					{
						CourierID:    1,
						OrderID:      1,
						CompleteTime: datetime.TimeAlias{},
					},
				},
			},
			assert.False,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.want(t, tc.req.Valid())
		})
	}
}

func TestCreateCourierRequest_Valid(t *testing.T) {
	tt := []struct {
		name string
		req  *CreateCourierRequest
		want assert.BoolAssertionFunc
	}{
		{"negative #1 - negative nil reference", nil, assert.False},
		{"negative #2 - negative no couriers", new(CreateCourierRequest), assert.False},
		{
			"positive #1 - one item",
			&CreateCourierRequest{
				Couriers: []CreateCourierDTO{
					{
						CourierType: FootCourierTypeString,
						Regions:     []int{1, 2},
						WorkingHours: []*datetime.TimeInterval{
							testTimeInterval1(t),
							testTimeInterval2(t),
						},
					},
				},
			},
			assert.True,
		},
		{
			"positive #2 - one item",
			&CreateCourierRequest{
				Couriers: []CreateCourierDTO{
					{
						CourierType: FootCourierTypeString,
						Regions:     []int{1, 2},
						WorkingHours: []*datetime.TimeInterval{
							testTimeInterval1(t),
							testTimeInterval2(t),
						},
					},
					{
						CourierType: AutoCourierTypeString,
						Regions:     []int{1, 2, 3},
						WorkingHours: []*datetime.TimeInterval{
							testTimeInterval1(t),
							testTimeInterval2(t),
						},
					},
				},
			},
			assert.True,
		},
		{
			"negative #3 - bad courier type",
			&CreateCourierRequest{
				Couriers: []CreateCourierDTO{
					{
						CourierType: "unknown type",
						Regions:     []int{1, 2},
						WorkingHours: []*datetime.TimeInterval{
							testTimeInterval1(t),
							testTimeInterval2(t),
						},
					},
					{
						CourierType: AutoCourierTypeString,
						Regions:     []int{1, 2, 3},
						WorkingHours: []*datetime.TimeInterval{
							testTimeInterval1(t),
							testTimeInterval2(t),
						},
					},
				},
			},
			assert.False,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.want(t, tc.req.Valid())
		})
	}
}

func TestCreateOrderDTO_Valid(t *testing.T) {
	tt := []struct {
		name  string
		order CreateOrderDTO
		want  assert.BoolAssertionFunc
	}{
		{
			"negative #1 - no WH",
			CreateOrderDTO{
				Weight:        0,
				Regions:       0,
				DeliveryHours: []*datetime.TimeInterval{},
				Cost:          0,
			},
			assert.False,
		},
		{
			"positive #1 - valid WH",
			CreateOrderDTO{
				Weight:  0,
				Regions: 0,
				DeliveryHours: []*datetime.TimeInterval{
					testTimeInterval1(t),
				},
				Cost: 0,
			},
			assert.True,
		},
		{
			"negative #2 - invalid WH",
			CreateOrderDTO{
				Weight:  0,
				Regions: 0,
				DeliveryHours: []*datetime.TimeInterval{
					testTimeInterval1(t),
					nil,
				},
				Cost: 0,
			},
			assert.False,
		},
		{
			"negative #3 - non distinct HW",
			CreateOrderDTO{
				Weight:  0,
				Regions: 0,
				DeliveryHours: []*datetime.TimeInterval{
					testTimeInterval1(t),
					testTimeInterval1(t),
				},
				Cost: 0,
			},
			assert.False,
		},
		{
			"negative #4 - negative weight",
			CreateOrderDTO{
				Weight:  -1,
				Regions: 0,
				DeliveryHours: []*datetime.TimeInterval{
					testTimeInterval1(t),
				},
				Cost: 0,
			},
			assert.False,
		},
		{
			"negative #5 - negative cost",
			CreateOrderDTO{
				Weight:  0,
				Regions: 0,
				DeliveryHours: []*datetime.TimeInterval{
					testTimeInterval1(t),
				},
				Cost: -1,
			},
			assert.False,
		},
		{
			"negative #5 - negative region",
			CreateOrderDTO{
				Weight:  0,
				Regions: -1,
				DeliveryHours: []*datetime.TimeInterval{
					testTimeInterval1(t),
				},
				Cost: 0,
			},
			assert.False,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.want(t, tc.order.Valid())
		})
	}
}

func TestCreateOrderRequest_Valid(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var req *CreateOrderRequest
		assert.False(t, req.Valid())
	})
	t.Run("empty req", func(t *testing.T) {
		req := new(CreateOrderRequest)
		assert.False(t, req.Valid())
	})
	t.Run("invalid order", func(t *testing.T) {
		req := &CreateOrderRequest{
			Orders: []CreateOrderDTO{
				{},
			},
		}
		assert.False(t, req.Valid())
	})
}
