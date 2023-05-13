package datetime

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	TestTimeInterval = &TimeInterval{
		start:   Minute(11*60 + 12),
		end:     Minute(22*60 + 33),
		reverse: false,
	}
	TestTimeTime1 = time.Date(2000, time.December, 1, 11, 12, 0, 0, time.UTC)
	TestTimeTime2 = time.Date(2000, time.December, 1, 22, 33, 0, 0, time.UTC)
)

func TestParseTimeInterval_OK_NonReversed(t *testing.T) {
	var startH, startM, endH, endM = 12, 59, 23, 33
	h, err := ParseTimeInterval(fmt.Sprintf("%d:%d-%d:%d", startH, startM, endH, endM))
	assert.NoError(t, err)
	if assert.NotNil(t, h) {
		assert.Equal(t, startH, h.start.Hour())
		assert.Equal(t, startM, h.start.Minute())
		assert.Equal(t, endH, h.end.Hour())
		assert.Equal(t, endM, h.end.Minute())
		assert.False(t, h.reverse)
	}
}

func TestParseTimeInterval_OK_Reversed(t *testing.T) {
	var endH, endM, startH, startM = 12, 59, 23, 33
	h, err := ParseTimeInterval(fmt.Sprintf("%d:%d-%d:%d", startH, startM, endH, endM))
	assert.NoError(t, err)
	if assert.NotNil(t, h) {
		assert.Equal(t, startH, h.start.Hour())
		assert.Equal(t, startM, h.start.Minute())
		assert.Equal(t, endH, h.end.Hour())
		assert.Equal(t, endM, h.end.Minute())
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

func TestTimeInterval_UnmarshalJSON_Positive(t *testing.T) {
	h := new(TimeInterval)

	err := h.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", TestTimeInterval)))
	assert.NoError(t, err)
	if assert.NotNil(t, h) {
		assert.Equal(t, TestTimeInterval, h)
	}
}

func TestTimeInterval_UnmarshalJSON_Negative_BadFormat(t *testing.T) {
	h := new(TimeInterval)

	err := h.UnmarshalJSON([]byte("{"))
	assert.Error(t, err)
}

func TestTimeInterval_UnmarshalJSON_Negative_BadTime(t *testing.T) {
	h := new(TimeInterval)

	err := h.UnmarshalJSON([]byte("\"12:22:22-22\""))
	assert.Error(t, err)
}

func TestTimeInterval_Start(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		assert.Equal(t, TestTimeInterval.start, TestTimeInterval.Start())
	})
	t.Run("nil time interval", func(t *testing.T) {
		assert.Equal(t, Minute(0), (*TimeInterval)(nil).Start())
	})
}

func TestTimeInterval_End(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		assert.Equal(t, TestTimeInterval.end, TestTimeInterval.End())
	})
	t.Run("nil time interval", func(t *testing.T) {
		assert.Equal(t, Minute(0), (*TimeInterval)(nil).End())
	})
}

func TestTimeIn(t *testing.T) {
	t.Run("main positive", func(t *testing.T) {
		assert.True(t, TestTimeInterval.TimeIn(TestTimeTime1))
		assert.True(t, TestTimeInterval.TimeIn(TestTimeTime2))
	})
}

func TestTimeIntervalAlias_TimeInterval(t *testing.T) {
	tic := TimeIntervalAlias{
		Start:   123,
		End:     321,
		Reverse: true,
	}
	ti := tic.TimeInterval()
	if assert.NotNil(t, ti) {
		assert.Equal(t, Minute(tic.End), ti.End())
		assert.Equal(t, Minute(tic.Start), ti.Start())
		assert.Equal(t, tic.Reverse, ti.reverse)
	}
}
