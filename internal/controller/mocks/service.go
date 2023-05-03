// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	middleware "github.com/vlad-marlo/yandex-academy-enrollment/internal/middleware"
	model "github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	datetime "github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/datetime"
)

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// BindAddr mocks base method.
func (m *MockConfig) BindAddr() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BindAddr")
	ret0, _ := ret[0].(string)
	return ret0
}

// BindAddr indicates an expected call of BindAddr.
func (mr *MockConfigMockRecorder) BindAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BindAddr", reflect.TypeOf((*MockConfig)(nil).BindAddr))
}

// MockServer is a mock of Server interface.
type MockServer struct {
	ctrl     *gomock.Controller
	recorder *MockServerMockRecorder
}

// MockServerMockRecorder is the mock recorder for MockServer.
type MockServerMockRecorder struct {
	mock *MockServer
}

// NewMockServer creates a new mock instance.
func NewMockServer(ctrl *gomock.Controller) *MockServer {
	mock := &MockServer{ctrl: ctrl}
	mock.recorder = &MockServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServer) EXPECT() *MockServerMockRecorder {
	return m.recorder
}

// Start mocks base method.
func (m *MockServer) Start(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockServerMockRecorder) Start(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockServer)(nil).Start), ctx)
}

// Stop mocks base method.
func (m *MockServer) Stop(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockServerMockRecorder) Stop(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockServer)(nil).Stop), ctx)
}

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AssignOrders mocks base method.
func (m *MockService) AssignOrders(ctx context.Context, date *datetime.Date) (*model.OrderAssignResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignOrders", ctx, date)
	ret0, _ := ret[0].(*model.OrderAssignResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignOrders indicates an expected call of AssignOrders.
func (mr *MockServiceMockRecorder) AssignOrders(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignOrders", reflect.TypeOf((*MockService)(nil).AssignOrders), ctx, date)
}

// CompleteOrders mocks base method.
func (m *MockService) CompleteOrders(ctx context.Context, req *model.CompleteOrderRequest) ([]*model.OrderDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteOrders", ctx, req)
	ret0, _ := ret[0].([]*model.OrderDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompleteOrders indicates an expected call of CompleteOrders.
func (mr *MockServiceMockRecorder) CompleteOrders(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteOrders", reflect.TypeOf((*MockService)(nil).CompleteOrders), ctx, req)
}

// CreateCouriers mocks base method.
func (m *MockService) CreateCouriers(ctx context.Context, request *model.CreateCourierRequest) (*model.CouriersCreateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCouriers", ctx, request)
	ret0, _ := ret[0].(*model.CouriersCreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCouriers indicates an expected call of CreateCouriers.
func (mr *MockServiceMockRecorder) CreateCouriers(ctx, request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCouriers", reflect.TypeOf((*MockService)(nil).CreateCouriers), ctx, request)
}

// CreateOrders mocks base method.
func (m *MockService) CreateOrders(ctx context.Context, req *model.CreateOrderRequest) ([]*model.OrderDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrders", ctx, req)
	ret0, _ := ret[0].([]*model.OrderDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrders indicates an expected call of CreateOrders.
func (mr *MockServiceMockRecorder) CreateOrders(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrders", reflect.TypeOf((*MockService)(nil).CreateOrders), ctx, req)
}

// GetCourierByID mocks base method.
func (m *MockService) GetCourierByID(ctx context.Context, id string) (*model.CourierDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCourierByID", ctx, id)
	ret0, _ := ret[0].(*model.CourierDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCourierByID indicates an expected call of GetCourierByID.
func (mr *MockServiceMockRecorder) GetCourierByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCourierByID", reflect.TypeOf((*MockService)(nil).GetCourierByID), ctx, id)
}

// GetCourierMetaInfo mocks base method.
func (m *MockService) GetCourierMetaInfo(ctx context.Context, req *model.GetCourierMetaInfoRequest) (*model.GetCourierMetaInfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCourierMetaInfo", ctx, req)
	ret0, _ := ret[0].(*model.GetCourierMetaInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCourierMetaInfo indicates an expected call of GetCourierMetaInfo.
func (mr *MockServiceMockRecorder) GetCourierMetaInfo(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCourierMetaInfo", reflect.TypeOf((*MockService)(nil).GetCourierMetaInfo), ctx, req)
}

// GetCouriers mocks base method.
func (m *MockService) GetCouriers(ctx context.Context, opts *middleware.PaginationOpts) (*model.GetCouriersResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCouriers", ctx, opts)
	ret0, _ := ret[0].(*model.GetCouriersResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCouriers indicates an expected call of GetCouriers.
func (mr *MockServiceMockRecorder) GetCouriers(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCouriers", reflect.TypeOf((*MockService)(nil).GetCouriers), ctx, opts)
}

// GetOrderByID mocks base method.
func (m *MockService) GetOrderByID(ctx context.Context, id string) (*model.OrderDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByID", ctx, id)
	ret0, _ := ret[0].(*model.OrderDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByID indicates an expected call of GetOrderByID.
func (mr *MockServiceMockRecorder) GetOrderByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByID", reflect.TypeOf((*MockService)(nil).GetOrderByID), ctx, id)
}

// GetOrders mocks base method.
func (m *MockService) GetOrders(ctx context.Context, opts *middleware.PaginationOpts) ([]*model.OrderDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", ctx, opts)
	ret0, _ := ret[0].([]*model.OrderDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockServiceMockRecorder) GetOrders(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockService)(nil).GetOrders), ctx, opts)
}

// GetOrdersAssign mocks base method.
func (m *MockService) GetOrdersAssign(ctx context.Context, date *datetime.Date, id string) (*model.OrderAssignResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersAssign", ctx, date, id)
	ret0, _ := ret[0].(*model.OrderAssignResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersAssign indicates an expected call of GetOrdersAssign.
func (mr *MockServiceMockRecorder) GetOrdersAssign(ctx, date, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersAssign", reflect.TypeOf((*MockService)(nil).GetOrdersAssign), ctx, date, id)
}
