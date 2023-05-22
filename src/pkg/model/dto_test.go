package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCourierDTO_EarningsConst(t *testing.T) {
	tt := []struct {
		name    string
		courier *CourierDTO
		want    int32
	}{
		{"nil courier", nil, unknownTypeConst},
		{"bad type", new(CourierDTO), unknownTypeConst},
		{"auto", &CourierDTO{CourierType: AutoCourierTypeString}, AutoCourierTypeEarningsConst},
		{"bike", &CourierDTO{CourierType: BikeCourierTypeString}, BikeCourierTypeEarningsConst},
		{"auto", &CourierDTO{CourierType: FootCourierTypeString}, FootCourierTypeEarningsConst},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.courier.EarningsConst())
		})
	}
}

func TestCourierDTO_RatingConst(t *testing.T) {
	tt := []struct {
		name    string
		courier *CourierDTO
		want    int32
	}{
		{"nil courier", nil, unknownTypeConst},
		{"bad type", new(CourierDTO), unknownTypeConst},
		{"auto", &CourierDTO{CourierType: AutoCourierTypeString}, AutoCourierTypeRatingConst},
		{"bike", &CourierDTO{CourierType: BikeCourierTypeString}, BikeCourierTypeRatingConst},
		{"auto", &CourierDTO{CourierType: FootCourierTypeString}, FootCourierTypeRatingConst},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.courier.RatingConst())
		})
	}
}
