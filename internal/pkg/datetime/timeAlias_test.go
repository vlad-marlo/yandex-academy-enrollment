package datetime

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTimeAlias_MarshalJSON(t *testing.T) {
	tt := []struct {
		name string
		time time.Time
	}{
		{"positive #1", time.Date(2000, time.December, 12, 22, 11, 34, 279, time.UTC)},
		{"positive #2", time.Date(2003, time.December, 12, 22, 23, 30, 0, time.UTC)},
		{"positive #3", time.Date(2014, time.December, 12, 23, 11, 30, 92, time.UTC)},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var got []byte

			jsonTime, err := json.Marshal(tc.time.Format(layout))
			require.NoError(t, err)
			got, err = (*TimeAlias)(&tc.time).MarshalJSON()
			require.NoError(t, err)
			assert.JSONEq(t, string(jsonTime), string(got))
		})
	}
}

func TestTimeAlias_MarshalJSON_Nil(t *testing.T) {
	b, err := (*TimeAlias)(nil).MarshalJSON()
	assert.Nil(t, b)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNilReference)
	}
}

func TestCorrectLayout(t *testing.T) {
	const str = "2023-05-08T09:39:12.260Z"
	parsed, err := time.Parse(layout, str)
	require.NoError(t, err)
	assert.Equal(t, parsed.Format(layout), str)
}

func TestTimeAlias_UnmarshalJSON_Positive(t *testing.T) {
	const str = "2023-05-08T09:39:12.260Z"

	myTime := new(TimeAlias)
	marshalled, err := json.Marshal(str)
	require.NoError(t, err)

	err = myTime.UnmarshalJSON(marshalled)
	require.NoError(t, err)
	assert.Equal(t, str, myTime.Time().Format(layout))
}

func TestTimeAlias_UnmarshalJSON_Negative_BadString(t *testing.T) {
	myTime := new(TimeAlias)
	assert.Error(t, myTime.UnmarshalJSON([]byte("some string")))
}

func TestTimeAlias_UnmarshalJSON_Negative_Nil(t *testing.T) {
	var myTime *TimeAlias
	const str = "2023-05-08T09:39:12.260Z"

	marshalled, err := json.Marshal(str)
	require.NoError(t, err)

	err = myTime.UnmarshalJSON(marshalled)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, ErrNilReference)
	}
	assert.Nil(t, myTime)
}

func TestTimeAlias_UnmarshalJSON_Negative(t *testing.T) {
	const str = "2023-05-08T09:39:12.260Z"
	tt := []struct {
		name string
		time *TimeAlias
		data []byte
		want assert.ErrorAssertionFunc
	}{
		{"nil time", nil, []byte(`"` + str + `"`), assert.Error},
		{"nil data", new(TimeAlias), nil, assert.Error},
		{"un-parsable time", new(TimeAlias), []byte("\"un-parsable time\""), assert.Error},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.time.UnmarshalJSON(tc.data)
			tc.want(t, err)
			t.Logf("got err: %v", err)
		})
	}
}

func TestTimeAlias_Time_Nil(t *testing.T) {
	var myTime *TimeAlias
	assert.Empty(t, myTime.Time())
}
