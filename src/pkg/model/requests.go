package model

import (
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
)

const (
	FootCourierTypeString = "FOOT"
	BikeCourierTypeString = "BIKE"
	AutoCourierTypeString = "AUTO"
)

type (
	CreateOrderRequest struct {
		Orders []CreateOrderDTO `json:"orders" validate:"required"`
	}
	CompleteOrderRequest struct {
		CompleteInfo []CompleteOrder `json:"complete_info" validate:"required"`
	}
	CompleteOrder struct {
		CourierID    int64         `json:"courier_id" validate:"required"`
		OrderID      int64         `json:"order_id" validate:"required"`
		CompleteTime datetime.Time `json:"complete_time,omitempty" swaggertype:"string" validate:"required"`
	}
	GetCourierRequest struct {
		CourierID int64 `path:"courier_id" validate:"required"`
	}
	GetCourierMetaInfoRequest struct {
		CourierID int64  `path:"courier_id" validate:"required"`
		StartDate string `query:"startDate" validate:"required"`
		EndDate   string `query:"endDate" validate:"required"`
	}
	CreateCourierRequest struct {
		Couriers []CreateCourierDTO `json:"couriers" validate:"required"`
	}
)
