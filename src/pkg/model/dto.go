package model

import (
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
)

type (
	OrderDTO struct {
		OrderID int64   `json:"order_id" validate:"required"`
		Weight  float64 `json:"weight" validate:"required"`
		Regions int32   `json:"regions" validate:"required"`
		// DeliveryHours is string slice of strings that represents time interval.
		//
		// String must be in HH:MM-HH:MM format where HH is hour (integer 0-23) and MM is minutes (integer 0-59).
		DeliveryHours []*datetime.TimeInterval `json:"delivery_hours" swaggertype:"array,string" validate:"required"`
		Cost          int32                    `json:"cost" validate:"required"`
		CompletedTime datetime.Time            `json:"completed_time,omitempty" swaggertype:"string"`
	}
	CreateOrderDTO struct {
		Weight  float64 `json:"weight" validate:"required"`
		Regions int32   `json:"regions" validate:"required"`
		// DeliveryHours is string slice of strings that represents time interval.
		//
		// String must be in HH:MM-HH:MM format where HH is hour (integer 0-23) and MM is minutes (integer 0-59).
		DeliveryHours []*datetime.TimeInterval `json:"delivery_hours" swaggertype:"array,string" validate:"required"`
		Cost          int32                    `json:"cost" validate:"required"`
	}

	CourierDTO struct {
		CourierID   int64   `json:"courier_id" validate:"required" example:"2"`
		CourierType string  `json:"courier_type" enums:"FOOT,BIKE,AUTO" validate:"required" example:"AUTO"`
		Regions     []int32 `json:"regions" validate:"required" example:"1,2,3"`
		// WorkingHours is string slice of strings that represents time interval.
		//
		// String must be in HH:MM-HH:MM format where HH is hour (integer 0-23) and MM is minutes (integer 0-59).
		WorkingHours []*datetime.TimeInterval `json:"working_hours" validate:"required" swaggertype:"array,string" example:"12:00-23:00,14:30-15:30"`
	}
	CreateCourierDTO struct {
		CourierType  string                   `json:"courier_type" enums:"FOOT,BIKE,AUTO" validate:"required" example:"AUTO"`
		Regions      []int32                  `json:"regions" validate:"required" example:"1,2,3"`
		WorkingHours []*datetime.TimeInterval `json:"working_hours" swaggertype:"array,string" validate:"required" example:"12:00-23:00,14:30-15:30"`
	}
)

const (
	unknownTypeConst             = 0
	FootCourierTypeEarningsConst = 1 + iota
	BikeCourierTypeEarningsConst
	AutoCourierTypeEarningsConst
	AutoCourierTypeRatingConst = iota - 3
	BikeCourierTypeRatingConst
	FootCourierTypeRatingConst
)

func (d *CourierDTO) EarningsConst() int32 {
	if d == nil {
		return unknownTypeConst
	}

	switch d.CourierType {
	case FootCourierTypeString:
		return FootCourierTypeEarningsConst
	case BikeCourierTypeString:
		return BikeCourierTypeEarningsConst
	case AutoCourierTypeString:
		return AutoCourierTypeEarningsConst
	default:
		return unknownTypeConst
	}
}

func (d *CourierDTO) RatingConst() int32 {
	if d == nil {
		return unknownTypeConst
	}

	switch d.CourierType {
	case AutoCourierTypeString:
		return AutoCourierTypeRatingConst
	case BikeCourierTypeString:
		return BikeCourierTypeRatingConst
	case FootCourierTypeString:
		return FootCourierTypeRatingConst
	default:
		return unknownTypeConst
	}
}
