package datetime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	TestTimeInterval = &TimeInterval{
		start: Time{
			hour:   11,
			minute: 12,
		},
		end: Time{
			hour:   22,
			minute: 33,
		},
		reverse: false,
	}
)

func TestParseTimeInterval_OK_NonReversed(t *testing.T) {
	var startH, startM, endH, endM uint8 = 12, 59, 23, 33
	h, err := ParseTimeInterval(fmt.Sprintf("%d:%d-%d:%d", startH, startM, endH, endM))
	assert.NoError(t, err)
	if assert.NotNil(t, h) {
		assert.Equal(t, startH, h.start.hour)
		assert.Equal(t, startM, h.start.minute)
		assert.Equal(t, endH, h.end.hour)
		assert.Equal(t, endM, h.end.minute)
		assert.False(t, h.reverse)
	}
}

func TestParseTimeInterval_OK_Reversed(t *testing.T) {
	var endH, endM, startH, startM uint8 = 12, 59, 23, 33
	h, err := ParseTimeInterval(fmt.Sprintf("%d:%d-%d:%d", startH, startM, endH, endM))
	assert.NoError(t, err)
	if assert.NotNil(t, h) {
		assert.Equal(t, startH, h.start.hour)
		assert.Equal(t, startM, h.start.minute)
		assert.Equal(t, endH, h.end.hour)
		assert.Equal(t, endM, h.end.minute)
		assert.True(t, h.reverse)
	}
}

func TestParseTimeInterval_BadData(t *testing.T) {
	t.Run("bad format", func(t *testing.T) {
		h, err := ParseTimeInterval("bad string")
		assert.Nil(t, h)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
	})
	t.Run("bad start time", func(t *testing.T) {
		h, err := ParseTimeInterval("bad time-12:23")
		assert.Nil(t, h)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
	})
	t.Run("bad end time", func(t *testing.T) {
		h, err := ParseTimeInterval("12:23-bad time")
		assert.Nil(t, h)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
	})
}

func TestTimeInterval_String_Parsable(t *testing.T) {
	h, err := ParseTimeInterval(TestTimeInterval.String())
	assert.NoError(t, err)
	if assert.NotNil(t, h) {
		assert.Equal(t, TestTimeInterval, h)
	}
}

func TestTimeInterval_MarshalJSON_Parsable(t *testing.T) {
	var h *TimeInterval

	data, err := TestTimeInterval.MarshalJSON()
	require.NoError(t, err)

	h, err = ParseTimeInterval(string(data))
	assert.NoError(t, err)
	if assert.NotNil(t, h) {
		assert.Equal(t, TestTimeInterval, h)
	}
}

func TestTimeInterval_UnmarshalJSON(t *testing.T) {
	h := new(TimeInterval)

	err := h.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", TestTimeInterval)))
	assert.NoError(t, err)
	if assert.NotNil(t, h) {
		assert.Equal(t, TestTimeInterval, h)
	}
}
