package production

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller/http"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/service/production/mocks"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/store"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/fielderr"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"math/rand"
	"testing"
)

var (
	testCourier = func(t testing.TB, id int64) *model.CourierDTO {
		t.Helper()
		return &model.CourierDTO{
			CourierID:   id,
			CourierType: model.FootCourierTypeString,
			Regions:     []int32{1, 3, 4, 144},
			WorkingHours: []*datetime.TimeInterval{
				datetime.TimeIntervalAlias{Start: 123, End: 321}.TimeInterval(),
				datetime.TimeIntervalAlias{Start: 403, End: 1233}.TimeInterval(),
			},
		}
	}
)

func couriersCreateToDTO(couriers []model.CreateCourierDTO) (res []model.CourierDTO) {
	for _, c := range couriers {
		res = append(res, model.CourierDTO{
			CourierID:    rand.Int63(),
			CourierType:  c.CourierType,
			Regions:      c.Regions,
			WorkingHours: c.WorkingHours,
		})
	}
	return
}

func TestService_GetCourierByID_Negative_NonParsableID(t *testing.T) {
	tt := []struct {
		name string
		id   string
	}{
		{"empty string", ""},
		{"non parsable string", "non parsable"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			srv := testService(t, nil)
			resp, err := srv.GetCourierByID(context.Background(), tc.id)
			if assert.Error(t, err) {
				assert.ErrorIs(t, err, ErrBadRequest)
			}
			assert.Nil(t, resp)
		})
	}
}

func TestService_GetCourierByID_Positive(t *testing.T) {
	type mockData struct {
		id      int64
		courier *model.CourierDTO
		err     error
	}
	type expect struct {
		courier *model.CourierDTO
		err     error
	}
	tt := []struct {
		name   string
		id     string
		mock   mockData
		expect expect
	}{
		{
			name:   "positive",
			id:     "123",
			mock:   mockData{123, testCourier(t, 123), nil},
			expect: expect{testCourier(t, 123), nil},
		},
		{
			name:   "negative",
			id:     "123",
			mock:   mockData{123, nil, ErrBadRequest},
			expect: expect{nil, ErrNotFound},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			str := mocks.NewMockStore(ctrl)
			str.EXPECT().GetCourierByID(gomock.Any(), tc.mock.id).Return(tc.mock.courier, tc.mock.err)

			srv := testService(t, str)

			resp, err := srv.GetCourierByID(context.Background(), tc.id)
			assert.Equal(t, tc.expect.courier, resp)
			assert.ErrorIs(t, err, tc.expect.err)
		})
	}
}

func TestService_CreateCouriers_Negative_NonValidRequest(t *testing.T) {
	tt := []struct {
		name string
		req  *model.CreateCourierRequest
	}{
		{"nil request", nil},
		{"non nil invalid request: nil couriers", &model.CreateCourierRequest{Couriers: nil}},
		{"non nil invalid request: zero size couriers", &model.CreateCourierRequest{Couriers: []model.CreateCourierDTO{}}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			srv := testService(t, nil)
			resp, err := srv.CreateCouriers(context.Background(), tc.req)
			assert.Nil(t, resp)
			if assert.Error(t, err) {
				assert.ErrorIs(t, err, ErrBadRequest)
			}
		})
	}
}

func TestService_CreateCouriers_Positive(t *testing.T) {
	tt := []struct {
		name     string
		couriers []model.CreateCourierDTO
	}{
		{
			name: "positive #1",
			couriers: []model.CreateCourierDTO{
				{
					CourierType: model.FootCourierTypeString,
					Regions:     []int32{3, 312},
					WorkingHours: []*datetime.TimeInterval{
						datetime.TimeIntervalAlias{Start: 12, End: 32}.TimeInterval(),
						datetime.TimeIntervalAlias{Start: 44, End: 123}.TimeInterval(),
					},
				},
			},
		},
		{
			name: "positive #2",
			couriers: []model.CreateCourierDTO{
				{
					CourierType: model.FootCourierTypeString,
					Regions:     []int32{3, 312},
					WorkingHours: []*datetime.TimeInterval{
						datetime.TimeIntervalAlias{Start: 12, End: 32}.TimeInterval(),
						datetime.TimeIntervalAlias{Start: 44, End: 123}.TimeInterval(),
					},
				},
				{
					CourierType: model.AutoCourierTypeString,
					Regions:     []int32{3, 312},
					WorkingHours: []*datetime.TimeInterval{
						datetime.TimeIntervalAlias{Start: 12, End: 32}.TimeInterval(),
						datetime.TimeIntervalAlias{Start: 44, End: 123}.TimeInterval(),
					},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			str := mocks.NewMockStore(ctrl)
			ctx := context.Background()
			couriers := couriersCreateToDTO(tc.couriers)
			str.EXPECT().CreateCouriers(ctx, tc.couriers).Return(couriers, nil)

			srv := testService(t, str)
			resp, err := srv.CreateCouriers(ctx, &model.CreateCourierRequest{Couriers: tc.couriers})
			assert.NoError(t, err)
			if assert.NotNil(t, resp) {
				assert.Equal(t, &model.CouriersCreateResponse{Couriers: couriers}, resp)
			}
		})
	}
}

func TestService_CreateCouriers_Negative_StorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	str.EXPECT().CreateCouriers(gomock.Any(), gomock.Any()).Return(nil, ErrNoContent)

	resp, err := srv.CreateCouriers(context.Background(), &model.CreateCourierRequest{Couriers: []model.CreateCourierDTO{{CourierType: model.FootCourierTypeString, Regions: []int32{3, 312}, WorkingHours: []*datetime.TimeInterval{datetime.TimeIntervalAlias{Start: 12, End: 32}.TimeInterval(), datetime.TimeIntervalAlias{Start: 44, End: 123}.TimeInterval()}}}})
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrBadRequest)
	}
	assert.Nil(t, resp)
}

func TestService_GetCouriers_NilOpts(t *testing.T) {
	srv := testService(t, nil)
	resp, err := srv.GetCouriers(context.Background(), nil)
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrBadRequest)
	}
}

func TestService_GetCouriers_StorageError_UnknownErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	str.EXPECT().GetCouriers(gomock.Any(), 1, 2).Return(nil, ErrNoContent)

	resp, err := srv.GetCouriers(context.Background(), http.NewPaginationOpts("1", "2"))
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrBadRequest)
	}
	assert.Nil(t, resp)
}

func TestService_GetCouriers_StorageError_NoContent(t *testing.T) {
	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	str.EXPECT().GetCouriers(gomock.Any(), 1, 2).Return(nil, store.ErrNoContent)

	want := &model.GetCouriersResponse{
		Couriers: []model.CourierDTO{},
		Limit:    1,
		Offset:   2,
	}

	resp, err := srv.GetCouriers(context.Background(), http.NewPaginationOpts("1", "2"))
	assert.NoError(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, want, resp)
	}
}

func TestService_GetCouriers_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	str.EXPECT().GetCouriers(gomock.Any(), 1, 2).Return([]model.CourierDTO{}, nil)

	want := &model.GetCouriersResponse{
		Couriers: []model.CourierDTO{},
		Limit:    1,
		Offset:   2,
	}

	resp, err := srv.GetCouriers(context.Background(), http.NewPaginationOpts("1", "2"))
	assert.NoError(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, want, resp)
	}
}

func TestService_GetCourierMetaInfo_Negative_NilReq(t *testing.T) {
	srv := testService(t, nil)
	resp, err := srv.GetCourierMetaInfo(context.Background(), nil)
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNoContent)
	}
}

func TestService_GetCourierMetaInfo_Negative_UnParsableStart(t *testing.T) {
	srv := testService(t, nil)
	resp, err := srv.GetCourierMetaInfo(context.Background(), &model.GetCourierMetaInfoRequest{})
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNoContent)
	}
}

func TestService_GetCourierMetaInfo_Negative_UnParsableEnd(t *testing.T) {
	srv := testService(t, nil)
	resp, err := srv.GetCourierMetaInfo(context.Background(), &model.GetCourierMetaInfoRequest{StartDate: "2022-12-31"})
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNoContent)
	}
}

func TestService_GetCourierMetaInfo_Negative_StorageErr_WhileGettingCourier(t *testing.T) {
	types := []string{model.BikeCourierTypeString, model.FootCourierTypeString, model.AutoCourierTypeString}
	for _, tp := range types {
		t.Run("check earnings and rating for type", func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			str := mocks.NewMockStore(ctrl)
			srv := testService(t, str)

			var resp *model.GetCourierMetaInfoResponse
			courier := model.CourierDTO{
				CourierID:    123,
				CourierType:  tp,
				Regions:      testCourier(t, 123).Regions,
				WorkingHours: testCourier(t, 123).WorkingHours,
			}
			req := &model.GetCourierMetaInfoRequest{
				CourierID: courier.CourierID,
				StartDate: "2022-12-31",
				EndDate:   "2023-01-01",
			}
			var end *datetime.Date
			start, err := datetime.ParseDate(req.StartDate)
			require.NoError(t, err)
			end, err = datetime.ParseDate(req.EndDate)
			require.NoError(t, err)

			str.EXPECT().GetCourierByID(gomock.Any(), courier.CourierID).Return(&courier, nil)
			str.EXPECT().GetCompletedOrdersPriceByCourier(gomock.Any(), courier.CourierID, start.Start(), end.End()).Return([]int32{100, 100, 100}, nil)

			resp, err = srv.GetCourierMetaInfo(ctx, req)
			assert.NoError(t, err)
			if assert.NotNil(t, resp) {
				earnings := 300 * courier.EarningsConst()
				rating := int32((3.0 / 48) * float64(courier.RatingConst()))
				want := &model.GetCourierMetaInfoResponse{
					CourierID:    courier.CourierID,
					CourierType:  courier.CourierType,
					Regions:      courier.Regions,
					WorkingHours: courier.WorkingHours,
					Rating:       rating,
					Earnings:     earnings,
				}
				assert.Equal(t, want, resp)
			}
		})
	}
}

func TestService_GetCourierMetaInfo_Positive_NoPrices(t *testing.T) {
	for _, returned := range []any{nil, []int32{}} {
		t.Run("with empty price", func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			str := mocks.NewMockStore(ctrl)
			srv := testService(t, str)

			var resp *model.GetCourierMetaInfoResponse
			courier := model.CourierDTO{
				CourierID:    123,
				CourierType:  model.BikeCourierTypeString,
				Regions:      testCourier(t, 123).Regions,
				WorkingHours: testCourier(t, 123).WorkingHours,
			}
			req := &model.GetCourierMetaInfoRequest{
				CourierID: courier.CourierID,
				StartDate: "2022-12-31",
				EndDate:   "2023-01-01",
			}
			var end *datetime.Date
			start, err := datetime.ParseDate(req.StartDate)
			require.NoError(t, err)
			end, err = datetime.ParseDate(req.EndDate)
			require.NoError(t, err)

			str.EXPECT().GetCourierByID(gomock.Any(), courier.CourierID).Return(&courier, nil)
			str.EXPECT().GetCompletedOrdersPriceByCourier(gomock.Any(), courier.CourierID, start.Start(), end.End()).Return(returned, nil)

			resp, err = srv.GetCourierMetaInfo(ctx, req)
			assert.NoError(t, err)
			if assert.NotNil(t, resp) {
				want := &model.GetCourierMetaInfoResponse{
					CourierID:    courier.CourierID,
					CourierType:  courier.CourierType,
					Regions:      courier.Regions,
					WorkingHours: courier.WorkingHours,
				}
				assert.Equal(t, want, resp)
			}
		})
	}
}

func TestService_GetCourierMetaInfo_ErrWhileGettingEarnings(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	var resp *model.GetCourierMetaInfoResponse
	courier := model.CourierDTO{
		CourierID:    123,
		CourierType:  model.BikeCourierTypeString,
		Regions:      testCourier(t, 123).Regions,
		WorkingHours: testCourier(t, 123).WorkingHours,
	}
	req := &model.GetCourierMetaInfoRequest{
		CourierID: courier.CourierID,
		StartDate: "2022-12-31",
		EndDate:   "2023-01-01",
	}
	var end *datetime.Date
	start, err := datetime.ParseDate(req.StartDate)
	require.NoError(t, err)
	end, err = datetime.ParseDate(req.EndDate)
	require.NoError(t, err)

	str.EXPECT().GetCourierByID(gomock.Any(), courier.CourierID).Return(&courier, nil)
	str.EXPECT().GetCompletedOrdersPriceByCourier(gomock.Any(), courier.CourierID, start.Start(), end.End()).Return(nil, store.ErrNoContent)

	resp, err = srv.GetCourierMetaInfo(ctx, req)
	if assert.Error(t, err) {
		if assert.ErrorIs(t, err, ErrNoContent) {
			var fErr *fielderr.Error
			if errors.As(err, &fErr) {
				_, ok := fErr.Data().(*model.GetCourierMetaInfoResponse)
				assert.True(t, ok)
			}
		}
	}
	assert.Nil(t, resp)

}

func TestService_GetCourierMetaInfo_ErrWhileGettingCourier(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	str := mocks.NewMockStore(ctrl)
	srv := testService(t, str)

	courier := model.CourierDTO{
		CourierID:    123,
		CourierType:  model.BikeCourierTypeString,
		Regions:      testCourier(t, 123).Regions,
		WorkingHours: testCourier(t, 123).WorkingHours,
	}
	req := &model.GetCourierMetaInfoRequest{
		CourierID: courier.CourierID,
		StartDate: "2022-12-31",
		EndDate:   "2023-01-01",
	}

	str.EXPECT().GetCourierByID(gomock.Any(), courier.CourierID).Return(nil, store.ErrNoContent)

	resp, err := srv.GetCourierMetaInfo(ctx, req)
	assert.Nil(t, resp)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNoContent)
	}
}
