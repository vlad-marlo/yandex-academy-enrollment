package model

const (
	FootCourierType = "FOOT"
	BikeCourierType = "BIKE"
	AutoCourierType = "AUTO"
)

type (
	CourierDTO struct {
		CourierID    int             `json:"courier_id"`
		CourierType  string          `json:"courier_type"`
		Regions      []int           `json:"regions"`
		WorkingHours []*TimeInterval `json:"working_hours"`
	}
	CourierGetRequest struct {
		CourierID int `path:"courier_id"`
	}
	CourierInCreateRequest struct {
		CourierType  string          `json:"courier_type"`
		Regions      []int           `json:"regions"`
		WorkingHours []*TimeInterval `json:"working_hours"`
	}
	CouriersCreateRequest struct {
		Couriers []CourierInCreateRequest `json:"couriers"`
	}
)

func (c *CourierInCreateRequest) regionsValid() bool {
	if c == nil {
		return false
	}
	return len(c.Regions) != 0
}

func (c *CourierInCreateRequest) workingHoursValid() bool {
	if c == nil {
		return false
	}
	return len(c.WorkingHours) != 0
}

func (c *CourierInCreateRequest) courierTypeValid() bool {
	switch c.CourierType {
	case FootCourierType, BikeCourierType, AutoCourierType:
	default:
		return false
	}
	return true
}

// Valid return is courier object valid.
func (c *CourierInCreateRequest) Valid() bool {
	return c.regionsValid() && c.courierTypeValid()
}
