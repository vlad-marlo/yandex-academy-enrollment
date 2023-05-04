package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/controller/mocks"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/datetime"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	TestTimeInterval1 *datetime.TimeInterval
	TestTimeInterval2 *datetime.TimeInterval
	TestCourier1      = &model.CourierDTO{
		CourierID:   1,
		CourierType: model.AutoCourierTypeString,
		Regions:     []int{1, 2, 3},
		WorkingHours: []*datetime.TimeInterval{
			TestTimeInterval1,
			TestTimeInterval2,
		},
	}
	TestDefaultPaginationParams = NewPaginationOpts("", "")
	ErrUnknown                  = errors.New("some unknown error")
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
		jsonCourier, err := json.Marshal(nil)
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
				Regions:      []int{1, 2, 3},
				WorkingHours: []*datetime.TimeInterval{},
			},
			{
				CourierID:    2,
				CourierType:  model.FootCourierTypeString,
				Regions:      []int{3, 4},
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
