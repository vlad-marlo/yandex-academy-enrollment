// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	model "github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CompleteOrders mocks base method.
func (m *MockStore) CompleteOrders(ctx context.Context, info []model.CompleteOrder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteOrders", ctx, info)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteOrders indicates an expected call of CompleteOrders.
func (mr *MockStoreMockRecorder) CompleteOrders(ctx, info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteOrders", reflect.TypeOf((*MockStore)(nil).CompleteOrders), ctx, info)
}

// CreateCouriers mocks base method.
func (m *MockStore) CreateCouriers(ctx context.Context, couriers []model.CreateCourierDTO) ([]model.CourierDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCouriers", ctx, couriers)
	ret0, _ := ret[0].([]model.CourierDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCouriers indicates an expected call of CreateCouriers.
func (mr *MockStoreMockRecorder) CreateCouriers(ctx, couriers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCouriers", reflect.TypeOf((*MockStore)(nil).CreateCouriers), ctx, couriers)
}

// CreateOrders mocks base method.
func (m *MockStore) CreateOrders(ctx context.Context, orders []*model.OrderDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrders", ctx, orders)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrders indicates an expected call of CreateOrders.
func (mr *MockStoreMockRecorder) CreateOrders(ctx, orders interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrders", reflect.TypeOf((*MockStore)(nil).CreateOrders), ctx, orders)
}

// GetCompletedOrdersPriceByCourier mocks base method.
func (m *MockStore) GetCompletedOrdersPriceByCourier(ctx context.Context, id int64, start, end time.Time) ([]int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompletedOrdersPriceByCourier", ctx, id, start, end)
	ret0, _ := ret[0].([]int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompletedOrdersPriceByCourier indicates an expected call of GetCompletedOrdersPriceByCourier.
func (mr *MockStoreMockRecorder) GetCompletedOrdersPriceByCourier(ctx, id, start, end interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompletedOrdersPriceByCourier", reflect.TypeOf((*MockStore)(nil).GetCompletedOrdersPriceByCourier), ctx, id, start, end)
}

// GetCourierByID mocks base method.
func (m *MockStore) GetCourierByID(ctx context.Context, id int64) (*model.CourierDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCourierByID", ctx, id)
	ret0, _ := ret[0].(*model.CourierDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCourierByID indicates an expected call of GetCourierByID.
func (mr *MockStoreMockRecorder) GetCourierByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCourierByID", reflect.TypeOf((*MockStore)(nil).GetCourierByID), ctx, id)
}

// GetCouriers mocks base method.
func (m *MockStore) GetCouriers(ctx context.Context, limit, offset int) ([]model.CourierDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCouriers", ctx, limit, offset)
	ret0, _ := ret[0].([]model.CourierDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCouriers indicates an expected call of GetCouriers.
func (mr *MockStoreMockRecorder) GetCouriers(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCouriers", reflect.TypeOf((*MockStore)(nil).GetCouriers), ctx, limit, offset)
}

// GetOrderByID mocks base method.
func (m *MockStore) GetOrderByID(ctx context.Context, id int64) (*model.OrderDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByID", ctx, id)
	ret0, _ := ret[0].(*model.OrderDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByID indicates an expected call of GetOrderByID.
func (mr *MockStoreMockRecorder) GetOrderByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByID", reflect.TypeOf((*MockStore)(nil).GetOrderByID), ctx, id)
}

// GetOrders mocks base method.
func (m *MockStore) GetOrders(ctx context.Context, limit, offset int) ([]*model.OrderDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", ctx, limit, offset)
	ret0, _ := ret[0].([]*model.OrderDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockStoreMockRecorder) GetOrders(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockStore)(nil).GetOrders), ctx, limit, offset)
}

// GetOrdersByIDs mocks base method.
func (m *MockStore) GetOrdersByIDs(ctx context.Context, ids []int64) ([]*model.OrderDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByIDs", ctx, ids)
	ret0, _ := ret[0].([]*model.OrderDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByIDs indicates an expected call of GetOrdersByIDs.
func (mr *MockStoreMockRecorder) GetOrdersByIDs(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByIDs", reflect.TypeOf((*MockStore)(nil).GetOrdersByIDs), ctx, ids)
}
