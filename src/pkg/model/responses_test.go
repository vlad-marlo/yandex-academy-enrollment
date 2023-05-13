package model

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrderDTO_ValidJSON(t *testing.T) {
	raw := fmt.Sprintf(
		`[
  {
    "order_id": 0,
    "weight": 0,
    "regions": 0,
    "delivery_hours": [
      "%s",
	  "%s"
    ],
    "cost": 0,
    "completed_time": "2023-05-08T16:51:45.198Z"
  }
]`,
		testTimeInterval1String,
		testTimeInterval2String,
	)

	var resp []OrderDTO
	require.NoError(t, json.Unmarshal([]byte(raw), &resp))
	got, err := json.Marshal(resp)
	require.NoError(t, err)
	assert.JSONEq(t, raw, string(got))
}
