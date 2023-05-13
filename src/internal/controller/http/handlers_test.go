package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller/mocks"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/fielderr"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	TestTimeInterval2 *datetime.TimeInterval

	TestTimeInterval1 = func(t testing.TB) *datetime.TimeInterval {
		h, err := datetime.ParseTimeInterval("12:22-13:30")
		require.NoError(t, err)
		return h
	}
	TestCourier1 = &model.CourierDTO{
		CourierID:   1,
		CourierType: model.AutoCourierTypeString,
		Regions:     []int32{1, 2, 3},
		WorkingHours: []*datetime.TimeInterval{
			TestTimeInterval2,
		},
	}
	TestOrder = func(t testing.TB) (string, *model.OrderDTO) {
		resp := &model.OrderDTO{
			OrderID: 123,
			Weight:  321,
			Regions: 412,
			DeliveryHours: []*datetime.TimeInterval{
				TestTimeInterval1(t),
			},
			Cost:          0,
			CompletedTime: datetime.Time{},
		}
		b, err := json.Marshal(resp)
		require.NoError(t, err)
		return string(b), resp
	}
	testOrders = func(t testing.TB) (string, []*model.OrderDTO) {
		_, order := TestOrder(t)
		resp := []*model.OrderDTO{order}
		b, err := json.Marshal(resp)
		require.NoError(t, err)
		return string(b), resp
	}
	testCreateOrderRequest = func(t testing.TB) (string, *model.CreateOrderRequest) {
		req := &model.CreateOrderRequest{
			Orders: []model.CreateOrderDTO{
				{
					Weight:  12,
					Regions: 312,
					DeliveryHours: []*datetime.TimeInterval{
						TestTimeInterval1(t),
					},
					Cost: 12,
				},
				{
					Weight:  3321,
					Regions: 312,
					DeliveryHours: []*datetime.TimeInterval{
						TestTimeInterval1(t),
					},
					Cost: 19023123,
				},
			},
		}
		b, err := json.Marshal(req)
		require.NoError(t, err)
		return string(b), req
	}
	TestDefaultPaginationParams = NewPaginationOpts("", "")
	ErrUnknown                  = errors.New("some unknown error")
	getOrdersAssignResponse     = func(t testing.TB) (string, *model.OrderAssignResponse) {
		t.Helper()
		resp := &model.OrderAssignResponse{}
		b, err := json.Marshal(resp)
		require.NoError(t, err)
		return string(b), resp
	}
	completeOrdersRequest = func(t testing.TB) (string, *model.CompleteOrderRequest) {
		req := &model.CompleteOrderRequest{
			CompleteInfo: []model.CompleteOrder{
				{1, 2, (datetime.Time)(time.Now())},
				{1, 21, (datetime.Time)(time.Now().Add(123 * time.Minute))},
				{231, 32213, (datetime.Time)(time.Now().Add(10 * time.Hour))},
			},
		}
		b, err := json.Marshal(req)
		require.NoError(t, err)
		return string(b), req
	}
	assignOrdersResponse = func(t testing.TB, date string) (string, *model.OrderAssignResponse) {
		resp := &model.OrderAssignResponse{
			Date: date,
			Couriers: []model.CourierGroupOrders{
				{
					CourierID: 2,
					Orders: []model.GroupOrders{
						{
							GroupOrderID: 1,
							Orders: []model.OrderDTO{
								{
									OrderID: 213,
									Weight:  1231412,
									Regions: 12414124,
									DeliveryHours: []*datetime.TimeInterval{
										TestTimeInterval1(t),
									},
									Cost:          1233,
									CompletedTime: (datetime.Time)(time.Now()),
								},
							},
						},
					},
				},
			},
		}
		b, err := json.Marshal(resp)
		require.NoError(t, err)
		return string(b), resp
	}
)

func TestController_HandleGetCourier_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	srv.EXPECT().GetCourierByID(gomock.Any(), fmt.Sprint(TestCourier1.CourierID)).Return(TestCourier1, nil)
	s := testServer(t, srv)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	defer assert.NoError(t, r.Body.Close())
	req := httptest.NewRecorder()
	c := s.engine.NewContext(r, req)
	c.SetParamNames("courier_id")
	c.SetParamValues(fmt.Sprint(TestCourier1.CourierID))
	if assert.NoError(t, s.HandleGetCourier(c)) {
		res := req.Result()
		defer assert.NoError(t, res.Body.Close())
		assert.Equal(t, http.StatusOK, req.Code)
		jsonCourier, err := json.Marshal(TestCourier1)
		require.NoError(t, err)
		assert.JSONEq(t, string(jsonCourier), req.Body.String())
	}
}

func TestController_HandleGetCourier_Negative_ErrInService(t *testing.T) {
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	srv.EXPECT().GetCourierByID(gomock.Any(), fmt.Sprint(TestCourier1.CourierID)).Return(nil, errors.New(""))
	s := testServer(t, srv)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	defer assert.NoError(t, r.Body.Close())
	req := httptest.NewRecorder()
	c := s.engine.NewContext(r, req)
	c.SetParamNames("courier_id")
	c.SetParamValues(fmt.Sprint(TestCourier1.CourierID))
	if assert.NoError(t, s.HandleGetCourier(c)) {
		res := req.Result()
		defer assert.NoError(t, res.Body.Close())
		assert.Equal(t, http.StatusNotFound, req.Code)
		jsonCourier, err := json.Marshal(model.BadRequestResponse{})
		require.NoError(t, err)
		assert.JSONEq(t, string(jsonCourier), req.Body.String())
	}
}

func TestController_HandlePing(t *testing.T) {
	srv := testServer(t, nil)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	defer assert.NoError(t, r.Body.Close())
	w := httptest.NewRecorder()
	c := srv.engine.NewContext(r, w)
	if assert.NoError(t, srv.HandlePing(c)) {
		res := w.Result()
		defer assert.NoError(t, res.Body.Close())
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestController_HandleGetCouriers_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	resp := &model.GetCouriersResponse{
		Couriers: []model.CourierDTO{
			{
				CourierID:    1,
				CourierType:  model.AutoCourierTypeString,
				Regions:      []int32{1, 2, 3},
				WorkingHours: []*datetime.TimeInterval{},
			},
			{
				CourierID:    2,
				CourierType:  model.FootCourierTypeString,
				Regions:      []int32{3, 4},
				WorkingHours: []*datetime.TimeInterval{},
			},
		},
		Limit:  1,
		Offset: 0,
	}
	srv.EXPECT().GetCouriers(gomock.Any(), gomock.Eq(TestDefaultPaginationParams)).Return(resp, nil)
	serv := testServer(t, srv)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	defer assert.NoError(t, r.Body.Close())
	c := serv.engine.NewContext(r, w)
	if assert.NoError(t, serv.HandleGetCouriers(c)) {
		res := w.Result()
		defer assert.NoError(t, res.Body.Close())
		wantBody, err := json.Marshal(resp)
		require.NoError(t, err)
		assert.JSONEq(t, string(wantBody), w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestController_HandleGetCouriers_Negative(t *testing.T) {
	// prepare mock storage
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	srv.EXPECT().GetCouriers(gomock.Any(), gomock.Eq(TestDefaultPaginationParams)).Return(nil, ErrUnknown)

	// prepare request
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	defer assert.NoError(t, r.Body.Close())

	// prepare server
	serv := testServer(t, srv)
	c := serv.engine.NewContext(r, w)
	err := serv.HandleGetCouriers(c)
	res := w.Result()
	defer assert.NoError(t, res.Body.Close())
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestController_HandleCreateCouriers_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	req := &model.CreateCourierRequest{
		Couriers: []model.CreateCourierDTO{
			{
				CourierType: model.AutoCourierTypeString,
				Regions:     []int32{1, 2, 3, 4},
				WorkingHours: []*datetime.TimeInterval{
					TestTimeInterval1(t),
				},
			},
			{
				CourierType: model.FootCourierTypeString,
				Regions:     []int32{3, 4, 12},
				WorkingHours: []*datetime.TimeInterval{
					TestTimeInterval1(t),
				},
			},
			{
				CourierType: model.BikeCourierTypeString,
				Regions:     []int32{7, 8},
				WorkingHours: []*datetime.TimeInterval{
					TestTimeInterval1(t),
				},
			},
		},
	}
	resp := &model.CouriersCreateResponse{
		Couriers: []model.CourierDTO{
			{
				CourierID:   1,
				CourierType: model.AutoCourierTypeString,
				Regions:     []int32{1, 2, 3, 4},
				WorkingHours: []*datetime.TimeInterval{
					TestTimeInterval1(t),
				},
			},
			{
				CourierID:   2,
				CourierType: model.FootCourierTypeString,
				Regions:     []int32{3, 4, 12},
				WorkingHours: []*datetime.TimeInterval{
					TestTimeInterval1(t),
				},
			},
			{
				CourierID:   3,
				CourierType: model.BikeCourierTypeString,
				Regions:     []int32{7, 8},
				WorkingHours: []*datetime.TimeInterval{
					TestTimeInterval1(t),
				},
			},
		},
	}
	srv.EXPECT().CreateCouriers(gomock.Any(), gomock.Eq(req)).Return(resp, nil)
	serv := testServer(t, srv)

	// prepare body
	body, err := json.Marshal(req)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	defer assert.NoError(t, r.Body.Close())
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())

	c := serv.engine.NewContext(r, w)
	if assert.NoError(t, serv.HandleCreateCouriers(c)) {
		res := w.Result()
		defer assert.NoError(t, res.Body.Close())
		assert.Equal(t, http.StatusOK, w.Code)
		wantBody, err := json.Marshal(resp)
		assert.NoError(t, err)
		assert.JSONEq(t, string(wantBody), w.Body.String(), w.Body.String())
	}
}

func TestController_HandleCreateCouriers_Negative_BadRequest(t *testing.T) {
	req := &model.CreateCourierRequest{
		Couriers: []model.CreateCourierDTO{
			{
				CourierType: model.AutoCourierTypeString,
				Regions:     []int32{1, 2, 3, 4},
				WorkingHours: []*datetime.TimeInterval{
					TestTimeInterval1(t),
				},
			},
			{
				CourierType: model.FootCourierTypeString,
				Regions:     []int32{3, 4, 12},
				WorkingHours: []*datetime.TimeInterval{
					TestTimeInterval1(t),
				},
			},
			{
				CourierType: model.BikeCourierTypeString,
				Regions:     []int32{7, 8},
				WorkingHours: []*datetime.TimeInterval{
					TestTimeInterval1(t),
				},
			},
		},
	}
	serv := testServer(t, nil)

	// prepare body
	body, err := json.Marshal(req)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	defer assert.NoError(t, r.Body.Close())
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())

	c := serv.engine.NewContext(r, w)
	if assert.NoError(t, serv.HandleCreateCouriers(c)) {
		res := w.Result()
		defer assert.NoError(t, res.Body.Close())
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestController_HandleCreateCouriers_Negative_ErrInService(t *testing.T) {
	tt := []struct {
		name       string
		err        error
		wantStatus int
		wantResp   interface{}
	}{
		{"unknown error", ErrUnknown, http.StatusBadRequest, model.BadRequestResponse{}},
		{"field error", fielderr.New("some error", someData, fielderr.CodeForbidden), http.StatusForbidden, someData},
		{"field error", fielderr.New("some error", nil, fielderr.CodeForbidden), http.StatusForbidden, nil},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// configure mock service
			ctrl := gomock.NewController(t)
			srv := mocks.NewMockService(ctrl)
			srv.EXPECT().CreateCouriers(gomock.Any(), gomock.Any()).Return(nil, tc.err)

			serv := testServer(t, srv)
			w := httptest.NewRecorder()
			defer assert.NoError(t, w.Result().Body.Close())
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			defer assert.NoError(t, r.Body.Close())
			c := serv.engine.NewContext(r, w)
			if assert.NoError(t, serv.HandleCreateCouriers(c)) {
				assert.Equal(t, tc.wantStatus, w.Code)
				wantResp, err := json.Marshal(tc.wantResp)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantResp), w.Body.String())
			}
		})
	}
}

func TestController_HandleGetCourierMetaInfo_Positive(t *testing.T) {
	var err error
	// prepare request and response objects
	var respString, reqString []byte
	req := &model.GetCourierMetaInfoRequest{
		CourierID: 1,
		StartDate: "2022-12-31",
		EndDate:   "2023-04-22",
	}
	reqString, err = json.Marshal(req)
	require.NoError(t, err)
	resp := &model.GetCourierMetaInfoResponse{
		CourierID:   1,
		CourierType: model.FootCourierTypeString,
		Regions:     []int32{1, 3, 2},
		WorkingHours: []*datetime.TimeInterval{
			TestTimeInterval1(t),
		},
		Rating:   12,
		Earnings: 1000,
	}
	respString, err = json.Marshal(resp)
	require.NoError(t, err)

	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	srv.EXPECT().GetCourierMetaInfo(gomock.Any(), gomock.Eq(req)).Return(resp, nil)

	serv := testServer(t, srv)

	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())

	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(reqString))
	defer assert.NoError(t, r.Body.Close())
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c := serv.engine.NewContext(r, w)
	if assert.NoError(t, serv.HandleGetCourierMetaInfo(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(respString), w.Body.String())
	}
}

func TestController_HandleGetCourierMetaInfo_Negative_BadRequest(t *testing.T) {
	var err error
	// prepare request and response objects
	var respString, reqString []byte
	req := &model.GetCourierMetaInfoRequest{
		CourierID: 1,
		StartDate: "2022-12-31",
		EndDate:   "2023-04-22",
	}
	reqString, err = json.Marshal(req)
	require.NoError(t, err)
	respString, err = json.Marshal(model.BadRequestResponse{})
	require.NoError(t, err)

	serv := testServer(t, nil)

	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())

	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(reqString))
	defer assert.NoError(t, r.Body.Close())

	c := serv.engine.NewContext(r, w)
	if assert.NoError(t, serv.HandleGetCourierMetaInfo(c)) {
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, string(respString), w.Body.String())
	}
}

func TestController_HandleGetCourierMetaInfo_Negative_InternalError(t *testing.T) {
	tt := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   interface{}
	}{
		{"unknown", ErrUnknown, http.StatusBadRequest, model.BadRequestResponse{}},
		{"fielderr", fielderr.New("some msg", someData, fielderr.CodeNotFound), http.StatusNotFound, someData},
		{"fielderr", fielderr.New("some msg", nil, fielderr.CodeNoContent), http.StatusNoContent, nil},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			// prepare request and response objects
			var respString, reqString []byte
			req := &model.GetCourierMetaInfoRequest{
				CourierID: 1,
				StartDate: "2022-12-31",
				EndDate:   "2023-04-22",
			}
			reqString, err = json.Marshal(req)
			require.NoError(t, err)
			respString, err = json.Marshal(tc.wantBody)
			require.NoError(t, err)

			ctrl := gomock.NewController(t)
			srv := mocks.NewMockService(ctrl)
			srv.EXPECT().GetCourierMetaInfo(gomock.Any(), gomock.Eq(req)).Return(nil, tc.err)

			serv := testServer(t, srv)

			w := httptest.NewRecorder()
			defer assert.NoError(t, w.Result().Body.Close())

			r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(reqString))
			defer assert.NoError(t, r.Body.Close())
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			c := serv.engine.NewContext(r, w)
			if assert.NoError(t, serv.HandleGetCourierMetaInfo(c)) {
				assert.Equal(t, tc.wantStatus, w.Code)
				assert.JSONEq(t, string(respString), w.Body.String())
			}
		})
	}
}

func TestController_HandleGetOrdersAssign_Positive_NoParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)

	wantResp, resp := getOrdersAssignResponse(t)
	srv.EXPECT().GetOrdersAssign(gomock.Any(), gomock.Eq(datetime.Today()), "").Return(resp, nil)

	serv := testServer(t, srv)
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	defer assert.NoError(t, r.Body.Close())
	c := serv.engine.NewContext(r, w)
	if assert.NoError(t, serv.HandleGetOrdersAssign(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, wantResp, w.Body.String())
	}
}

func TestController_HandleGetOrdersAssign_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)

	wantResp, resp := getOrdersAssignResponse(t)
	date, err := datetime.ParseDate("2022-12-30")
	require.NoError(t, err)
	srv.EXPECT().GetOrdersAssign(gomock.Any(), gomock.Eq(date), "1").Return(resp, nil)

	serv := testServer(t, srv)
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())
	r := httptest.NewRequest(http.MethodGet, "/?courier_id=1&date=2022-12-30", nil)
	defer assert.NoError(t, r.Body.Close())
	c := serv.engine.NewContext(r, w)
	if assert.NoError(t, serv.HandleGetOrdersAssign(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, wantResp, w.Body.String())
	}
}

func TestController_HandleGetOrdersAssign_NegativeBadDate(t *testing.T) {
	serv := testServer(t, nil)
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())
	r := httptest.NewRequest(http.MethodGet, "/?date=2022-02-29", nil)
	defer assert.NoError(t, r.Body.Close())
	c := serv.engine.NewContext(r, w)
	if assert.NoError(t, serv.HandleGetOrdersAssign(c)) {
		assert.Equal(t, http.StatusBadRequest, w.Code)
	}
}

func TestController_HandleGetOrdersAssign_Negative_InternalErr(t *testing.T) {
	tt := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   interface{}
	}{
		{"unknown", ErrUnknown, http.StatusBadRequest, model.BadRequestResponse{}},
		{"fielderr", fielderr.New("some msg", someData, fielderr.CodeNotFound), http.StatusNotFound, someData},
		{"fielderr", fielderr.New("some msg", nil, fielderr.CodeNoContent), http.StatusNoContent, nil},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			const (
				dateLayout = "2022-02-28"
				courierID  = "123"
			)

			ctrl := gomock.NewController(t)
			srv := mocks.NewMockService(ctrl)

			date, err := datetime.ParseDate(dateLayout)
			require.NoError(t, err)
			srv.EXPECT().GetOrdersAssign(gomock.Any(), gomock.Eq(date), courierID).Return(nil, tc.err)

			serv := testServer(t, srv)
			w := httptest.NewRecorder()
			defer assert.NoError(t, w.Result().Body.Close())
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?courier_id=%s&date=%s", courierID, dateLayout), nil)
			defer assert.NoError(t, r.Body.Close())
			c := serv.engine.NewContext(r, w)
			if assert.NoError(t, serv.HandleGetOrdersAssign(c)) {
				assert.Equal(t, tc.wantStatus, w.Code)
				wantBody, err := json.Marshal(tc.wantBody)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantBody), w.Body.String())
			}
		})
	}

}

func TestController_HandleGetOrder_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	const id = "123"
	wantResp, resp := TestOrder(t)
	srv.EXPECT().GetOrderByID(gomock.Any(), id).Return(resp, nil)

	serv := testServer(t, srv)

	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	defer assert.NoError(t, r.Body.Close())
	c := serv.engine.NewContext(r, w)
	c.SetParamNames("order_id")
	c.SetParamValues(id)
	if assert.NoError(t, serv.HandleGetOrder(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, wantResp, w.Body.String())
	}
}

func TestController_HandleGetOrder_Negative_NF(t *testing.T) {
	tt := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   interface{}
	}{
		{"unknown", ErrUnknown, http.StatusBadRequest, model.BadRequestResponse{}},
		{"fielderr", fielderr.New("some msg", someData, fielderr.CodeNotFound), http.StatusNotFound, someData},
		{"fielderr", fielderr.New("some msg", nil, fielderr.CodeNoContent), http.StatusNoContent, nil},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			srv := mocks.NewMockService(ctrl)
			const id = "123"
			srv.EXPECT().GetOrderByID(gomock.Any(), id).Return(nil, tc.err)

			serv := testServer(t, srv)

			w := httptest.NewRecorder()
			defer assert.NoError(t, w.Result().Body.Close())
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			defer assert.NoError(t, r.Body.Close())
			c := serv.engine.NewContext(r, w)
			c.SetParamNames("order_id")
			c.SetParamValues(id)
			if assert.NoError(t, serv.HandleGetOrder(c)) {
				assert.Equal(t, tc.wantStatus, w.Code)
				wantBody, err := json.Marshal(tc.wantBody)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantBody), w.Body.String())
			}
		})
	}
}

func TestController_HandleGetOrders_Positive(t *testing.T) {
	wantResp, resp := testOrders(t)
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	srv.EXPECT().GetOrders(gomock.Any(), gomock.Eq(NewPaginationOpts("", ""))).Return(resp, nil)

	serv := testServer(t, srv)

	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	defer assert.NoError(t, r.Body.Close())

	c := serv.engine.NewContext(r, w)
	if assert.NoError(t, serv.HandleGetOrders(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, wantResp, w.Body.String())
	}
}

func TestController_HandleGetOrders_Negative(t *testing.T) {
	tt := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   interface{}
	}{
		{"unknown", ErrUnknown, http.StatusBadRequest, model.BadRequestResponse{}},
		{"fielderr", fielderr.New("some msg", someData, fielderr.CodeNotFound), http.StatusNotFound, someData},
		{"fielderr", fielderr.New("some msg", nil, fielderr.CodeNoContent), http.StatusNoContent, nil},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			srv := mocks.NewMockService(ctrl)
			srv.EXPECT().GetOrders(gomock.Any(), gomock.Eq(NewPaginationOpts("", ""))).Return(nil, tc.err)

			serv := testServer(t, srv)

			w := httptest.NewRecorder()
			defer assert.NoError(t, w.Result().Body.Close())
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			defer assert.NoError(t, r.Body.Close())

			c := serv.engine.NewContext(r, w)
			if assert.NoError(t, serv.HandleGetOrders(c)) {
				assert.Equal(t, tc.wantStatus, w.Code)
				wantResp, err := json.Marshal(tc.wantBody)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantResp), w.Body.String())
			}
		})
	}
}

func TestController_HandleCreateOrders_Positive(t *testing.T) {
	body, req := testCreateOrderRequest(t)
	wantResp, resp := testOrders(t)

	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	srv.EXPECT().CreateOrders(gomock.Any(), gomock.Eq(req)).Return(resp, nil)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	defer assert.NoError(t, r.Body.Close())
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())

	serv := testServer(t, srv)
	c := serv.engine.NewContext(r, w)

	if assert.NoError(t, serv.HandleCreateOrders(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, wantResp, w.Body.String())
	}
}

func TestController_HandleCreateOrders_Negative_BadRequest(t *testing.T) {
	body, _ := testCreateOrderRequest(t)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	defer assert.NoError(t, r.Body.Close())
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())

	serv := testServer(t, nil)
	c := serv.engine.NewContext(r, w)

	if assert.NoError(t, serv.HandleCreateOrders(c)) {
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, "{}", w.Body.String())
	}
}

func TestController_HandleCreateOrders_Negative_Internal(t *testing.T) {
	tt := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   interface{}
	}{
		{"unknown", ErrUnknown, http.StatusBadRequest, model.BadRequestResponse{}},
		{"fielderr", fielderr.New("some msg", someData, fielderr.CodeNotFound), http.StatusNotFound, someData},
		{"fielderr", fielderr.New("some msg", nil, fielderr.CodeNoContent), http.StatusNoContent, nil},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			body, req := testCreateOrderRequest(t)

			ctrl := gomock.NewController(t)
			srv := mocks.NewMockService(ctrl)
			srv.EXPECT().CreateOrders(gomock.Any(), gomock.Eq(req)).Return(nil, tc.err)

			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			defer assert.NoError(t, r.Body.Close())
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			defer assert.NoError(t, w.Result().Body.Close())

			serv := testServer(t, srv)
			c := serv.engine.NewContext(r, w)

			if assert.NoError(t, serv.HandleCreateOrders(c)) {
				assert.Equal(t, tc.wantStatus, w.Code)
				wantResp, err := json.Marshal(tc.wantBody)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantResp), w.Body.String())
			}
		})
	}
}

func TestController_HandleCompleteOrders_Positive(t *testing.T) {
	body, _ := completeOrdersRequest(t)
	wantResp, resp := testOrders(t)
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	srv.EXPECT().CompleteOrders(gomock.Any(), gomock.Any()).Return(resp, nil)

	serv := testServer(t, srv)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	defer assert.NoError(t, r.Body.Close())
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())
	c := serv.engine.NewContext(r, w)

	if assert.NoError(t, serv.HandleCompleteOrders(c)) {
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, wantResp, w.Body.String())
	}
}

func TestController_HandleCompleteOrders_Negative_BadRequest(t *testing.T) {
	body, _ := completeOrdersRequest(t)

	serv := testServer(t, nil)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	defer assert.NoError(t, r.Body.Close())
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())
	c := serv.engine.NewContext(r, w)

	if assert.NoError(t, serv.HandleCompleteOrders(c)) {
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, "{}", w.Body.String())
	}
}

func TestController_HandleCompleteOrders_Negative_ErrInternal(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		wantCode int
		wantBody interface{}
	}{
		{"unknown", ErrUnknown, http.StatusBadRequest, model.BadRequestResponse{}},
		{"fielderr", fielderr.New("some msg", someData, fielderr.CodeNotFound), http.StatusNotFound, someData},
		{"fielderr", fielderr.New("some msg", nil, fielderr.CodeNoContent), http.StatusNoContent, nil},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := completeOrdersRequest(t)
			ctrl := gomock.NewController(t)
			srv := mocks.NewMockService(ctrl)
			srv.EXPECT().CompleteOrders(gomock.Any(), gomock.Any()).Return(nil, tc.err)

			serv := testServer(t, srv)

			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			defer assert.NoError(t, r.Body.Close())
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			defer assert.NoError(t, w.Result().Body.Close())
			c := serv.engine.NewContext(r, w)

			if assert.NoError(t, serv.HandleCompleteOrders(c)) {
				assert.Equal(t, tc.wantCode, w.Code)
				wantResp, err := json.Marshal(tc.wantBody)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantResp), w.Body.String())
			}
		})
	}
}

func TestController_HandleAssignOrders_Positive(t *testing.T) {
	const dateLayout = "2022-12-22"
	date, err := datetime.ParseDate(dateLayout)
	require.NoError(t, err)

	wantResp, resp := assignOrdersResponse(t, dateLayout)

	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)

	srv.EXPECT().AssignOrders(gomock.Any(), gomock.Eq(date)).Return(resp, nil)

	r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/?date=%s", dateLayout), nil)
	defer assert.NoError(t, r.Body.Close())
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())

	serv := testServer(t, srv)
	c := serv.engine.NewContext(r, w)

	if assert.NoError(t, serv.HandleAssignOrders(c)) {
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, wantResp, w.Body.String())
	}
}

func TestController_HandleAssignOrders_Negative_BadDate(t *testing.T) {
	const dateLayout = "2022-12-33"

	r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/?date=%s", dateLayout), nil)
	defer assert.NoError(t, r.Body.Close())
	w := httptest.NewRecorder()
	defer assert.NoError(t, w.Result().Body.Close())

	serv := testServer(t, nil)
	c := serv.engine.NewContext(r, w)

	if assert.NoError(t, serv.HandleAssignOrders(c)) {
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, "{}", w.Body.String())
	}
}

func TestController_HandleAssignOrders_Negative_ErrInternal(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		wantCode int
		wantBody interface{}
	}{
		{"unknown", ErrUnknown, http.StatusBadRequest, model.BadRequestResponse{}},
		{"fielderr", fielderr.New("some msg", someData, fielderr.CodeNotFound), http.StatusNotFound, someData},
		{"fielderr", fielderr.New("some msg", nil, fielderr.CodeNoContent), http.StatusNoContent, nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			const dateLayout = "2022-12-22"
			date, err := datetime.ParseDate(dateLayout)
			require.NoError(t, err)

			ctrl := gomock.NewController(t)
			srv := mocks.NewMockService(ctrl)

			srv.EXPECT().AssignOrders(gomock.Any(), gomock.Eq(date)).Return(nil, tc.err)

			r := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/?date=%s", dateLayout), nil)
			defer assert.NoError(t, r.Body.Close())
			w := httptest.NewRecorder()
			defer assert.NoError(t, w.Result().Body.Close())

			serv := testServer(t, srv)
			c := serv.engine.NewContext(r, w)

			if assert.NoError(t, serv.HandleAssignOrders(c)) {
				assert.Equal(t, tc.wantCode, w.Code)
				wantResp, err := json.Marshal(tc.wantBody)
				require.NoError(t, err)
				assert.JSONEq(t, string(wantResp), w.Body.String())
			}
		})
	}
}
