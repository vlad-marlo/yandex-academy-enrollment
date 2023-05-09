package model

import "github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/datetime"

type (
	GroupOrders struct {
		GroupOrderID int        `json:"group_order_id"`
		Orders       []OrderDTO `json:"orders" validate:"required"`
	}
	CourierGroupOrders struct {
		CourierID int           `json:"courier_id"`
		Orders    []GroupOrders `json:"orders"`
	}
	CouriersCreateResponse struct {
		Couriers []CourierDTO `json:"couriers"`
	}
	BadRequestResponse  struct{}
	GetCouriersResponse struct {
		Couriers []CourierDTO `json:"couriers"`
		Limit    int          `json:"limit"`
		Offset   int          `json:"offset"`
	}
	GetCourierMetaInfoResponse struct {
		CourierID   int    `json:"courier_id" validate:"required" example:"1"`
		CourierType string `json:"courier_type" enums:"FOOT,BIKE,AUTO" validate:"required" example:"AUTO"`
		Regions     []int  `json:"regions" validate:"required" example:"1,3,6"`
		// WorkingHours is string slice of strings that represents time interval.
		//
		// String must be in HH:MM-HH:MM format where HH is hour (integer 0-23) and MM is minutes (integer 0-59).
		WorkingHours []*datetime.TimeInterval `json:"working_hours" validate:"required" swaggertype:"array,string" example:"12:00-23:00,14:30-15:30"`
		Rating       int                      `json:"rating,omitempty"`
		Earnings     int                      `json:"earnings,omitempty"`
	}
	OrderAssignResponse struct {
		Date     string               `json:"date"`
		Couriers []CourierGroupOrders `json:"couriers"`
	}
)
