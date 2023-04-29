package model

const (
	FootCourierType = "FOOT"
	BikeCourierType = "BIKE"
	AutoCourierType = "AUTO"
)

type (
	Courier struct {
		CourierID    int             `json:"courier_id"`
		CourierType  string          `json:"courier_type"`
		Regions      []int           `json:"regions"`
		WorkingHours []*TimeInterval `json:"working_hours"`
	}
	CouriersCreateRequest []Courier
)
