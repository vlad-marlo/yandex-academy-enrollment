package model

type Courier struct {
	CourierID    int      `json:"courier_id"`
	CourierType  string   `json:"courier_type"`
	Regions      []int    `json:"regions"`
	WorkingHours []string `json:"working_hours"`
}
