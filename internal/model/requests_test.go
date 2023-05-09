package model

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/datetime"
	"testing"
)

func TestCreateOrderRequest_Marshalling(t *testing.T) {
	raw := fmt.Sprintf(`{
  "orders": [
    {
      "weight": 11,
      "regions": 22,
      "delivery_hours": [
        "%s",
		"%s"
      ],
      "cost": 123
    }
  ]
}`, testTimeInterval1String, testTimeInterval2String)
	want := &CreateOrderRequest{
		Orders: []CreateOrderDTO{
			{
				Weight:  11,
				Regions: 22,
				DeliveryHours: []*datetime.TimeInterval{
					testTimeInterval1(t),
					testTimeInterval2(t),
				},
				Cost: 123,
			},
		},
	}
	req := new(CreateOrderRequest)
	require.NoError(t, json.Unmarshal([]byte(raw), req))
	assert.Equal(t, want, req)
	assert.True(t, req.Valid())
}
