package model

import (
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/collections"
)

var typeSet = collections.NewSet[string](FootCourierTypeString, AutoCourierTypeString, BikeCourierTypeString)

func all(f func(item any) bool, iter []any) bool {
	for _, i := range iter {
		if !f(i) {
			return false
		}
	}
	return true
}

// Valid validates request.
//
// It is nilness safe function.
func (req *CompleteOrderRequest) Valid() bool {
	if req == nil {
		return false
	}
	l := len(req.CompleteInfo)
	switch l {
	case 0:
		return false
	case 1:
		return true
	default:
	}
	couriers := collections.NewSet[int]()
	orders := collections.NewSet[int]()
	for _, i := range req.CompleteInfo {
		if orders.Contain(i.OrderID) || couriers.Contain(i.CourierID) {
			return false
		}
		orders.Add(i.OrderID)
		couriers.Add(i.CourierID)
	}
	return true
}

// Valid validates request.
//
// It is nilness safe function.
func (d CreateCourierDTO) Valid() bool {
	return collections.Distinct[int](d.Regions...) && len(d.Regions) > 0 && typeSet.Contain(d.CourierType)
}

// Valid validates request.
//
// It is nilness safe function.
func (req *CreateCourierRequest) Valid() bool {
	if req == nil {
		return false
	}
	for _, i := range req.Couriers {
		if !i.Valid() {
			return false
		}
	}
	return len(req.Couriers) > 0
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
	for _, o := range req.Orders {
		if !o.Valid() {
			return false
		}
	}
	return len(req.Orders) > 0
}
