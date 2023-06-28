package model

import (
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/collections"
	"golang.org/x/exp/constraints"
)

var typeSet = collections.NewSet[string](FootCourierTypeString, AutoCourierTypeString, BikeCourierTypeString)

type anyFromModels interface {
	CreateOrderDTO | OrderDTO | CourierDTO | CreateCourierDTO |
		CreateOrderRequest | CompleteOrderRequest | CreateCourierRequest | CompleteOrder |
		GetCourierMetaInfoRequest
}

func all[T anyFromModels | constraints.Integer | constraints.Float](f func(item T) bool, iter ...T) bool {
	for _, i := range iter {
		if !f(i) {
			return false
		}
	}
	return len(iter) > 0
}

// Valid validates request.
//
// It is nilness safe function.
func (req *CompleteOrderRequest) Valid() bool {
	if req == nil {
		return false
	}
	orders := collections.NewSet[int64]()
	return all[CompleteOrder](
		func(item CompleteOrder) bool {
			defer orders.Add(item.OrderID)
			return !orders.Contain(item.OrderID)
		},
		req.CompleteInfo...,
	)
}

// Valid validates request.
//
// It is nilness safe function.
func (d CreateCourierDTO) Valid() bool {
	return collections.Distinct[int32](d.Regions...) && len(d.Regions) > 0 && typeSet.Contain(d.CourierType)
}

// Valid validates request.
//
// It is nilness safe function.
func (req *CreateCourierRequest) Valid() bool {
	if req == nil {
		return false
	}
	return all[CreateCourierDTO](
		func(item CreateCourierDTO) bool { return item.Valid() },
		req.Couriers...,
	)
}

// Valid validates request.
//
// It is nilness safe function.
func (d CreateOrderDTO) Valid() bool {
	set := collections.NewSet[string]()
	for _, h := range d.DeliveryHours {
		if h == nil {
			return false
		}
		set.Add(h.String())
	}
	ok := set.Len() == len(d.DeliveryHours) && len(d.DeliveryHours) > 0
	ok = ok && d.Weight >= 0 && d.Regions >= 0 && d.Cost >= 0
	return ok
}

// Valid validates request.
//
// It is nilness safe function.
func (req *CreateOrderRequest) Valid() bool {
	if req == nil {
		return false
	}
	return all[CreateOrderDTO](
		func(item CreateOrderDTO) bool { return item.Valid() },
		req.Orders...,
	)
}
