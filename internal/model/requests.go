package model

import (
	"time"
)

const (
	FootCourierType = iota
	BikeCourierType
	AutoCourierType
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
		CourierID    int       `json:"courier_id" validate:"required"`
		OrderID      int       `json:"order_id" validate:"required"`
		CompleteTime time.Time `json:"complete_time" validate:"required"`
	}

	GetCourierRequest struct {
		CourierID int `path:"courier_id" validate:"required"`
	}

	GetCourierMetaInfoRequest struct {
		CourierID int    `path:"courier_id" validate:"required"`
		StartDate string `query:"startDate" validate:"required"`
		EndDate   string `query:"endDate" validate:"required"`
	}
	CreateCourierRequest struct {
		Couriers []CreateCourierDTO `json:"couriers" validate:"required"`
	}
)
