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
)

func TestController_HandleGetCourier_Positive(t *testing.T) {
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	srv.EXPECT().GetCourierByID(gomock.Any(), TestCourier1.CourierID).Return(TestCourier1, nil)
	s := testServer(t, srv)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httptest.NewRecorder()
	c := s.engine.NewContext(r, req)
	c.SetParamNames("courier_id")
	c.SetParamValues(fmt.Sprint(TestCourier1.CourierID))
	if assert.NoError(t, s.HandleGetCourier(c)) {
		assert.Equal(t, http.StatusOK, req.Code)
		jsonCourier, err := json.Marshal(TestCourier1)
		require.NoError(t, err)
		assert.JSONEq(t, string(jsonCourier), req.Body.String())
	}
}

func TestController_HandleGetCourier_Negative_ErrInService(t *testing.T) {
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	srv.EXPECT().GetCourierByID(gomock.Any(), TestCourier1.CourierID).Return(nil, errors.New(""))
	s := testServer(t, srv)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httptest.NewRecorder()
	c := s.engine.NewContext(r, req)
	c.SetParamNames("courier_id")
	c.SetParamValues(fmt.Sprint(TestCourier1.CourierID))
	if assert.NoError(t, s.HandleGetCourier(c)) {
		assert.Equal(t, http.StatusNotFound, req.Code)
		jsonCourier, err := json.Marshal(nil)
		require.NoError(t, err)
		assert.JSONEq(t, string(jsonCourier), req.Body.String())
	}
}

func TestController_HandleGetCourier_Negative_ErrWhileParsingInt(t *testing.T) {
	ctrl := gomock.NewController(t)
	srv := mocks.NewMockService(ctrl)
	s := testServer(t, srv)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httptest.NewRecorder()
	c := s.engine.NewContext(r, req)
	if assert.NoError(t, s.HandleGetCourier(c)) {
		assert.Equal(t, http.StatusNotFound, req.Code)
		jsonCourier, err := json.Marshal(nil)
		require.NoError(t, err)
		assert.JSONEq(t, string(jsonCourier), req.Body.String())
	}
}
