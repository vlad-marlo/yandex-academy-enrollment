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

func TestMin(t *testing.T) {
	for x := 0; x != 1000; x++ {
		for y := 0; y != 1000; y++ {
			if x > y {
				assert.Equal(t, y, min(x, y))
				assert.Equal(t, x, max(x, y))
			} else {
				assert.Equal(t, x, min(x, y))
				assert.Equal(t, y, max(x, y))
			}
		}
	}
}

func TestTimeInterval_Common_Positive(t *testing.T) {
	var tt = []struct {
		name    string
		h       *TimeInterval
		other   *TimeInterval
		want    *TimeInterval
		minutes Minute
	}{
		{
			name: "non-reversed intervals with common time interval",
			h: &TimeInterval{
				start: 12 * 60,
				end:   14*60 + 30,
			},
			other: &TimeInterval{
				start: 13*60 + 30,
				end:   15*60 + 30,
			},
			want: &TimeInterval{
				start:   13*60 + 30,
				end:     14*60 + 30,
				reverse: false,
			},
			minutes: 61,
		},
		{
			name: "no common interval",
			h: &TimeInterval{
				start: 0,
				end:   1,
			},
			other: &TimeInterval{
				start: 2,
				end:   4,
			},
			want:    new(TimeInterval),
			minutes: 0,
		},
		{
			name: "only one common minute",
			h: &TimeInterval{
				start: 0,
				end:   1,
			},
			other: &TimeInterval{
				start: 1,
				end:   2,
			},
			want: &TimeInterval{
				start: 1,
				end:   1,
			},
			minutes: 1,
		},
		{
			name: "two reversed with collision",
			h: &TimeInterval{
				start:   1,
				end:     0,
				reverse: true,
			},
			other: &TimeInterval{
				start:   2,
				end:     1,
				reverse: true,
			},
			want: &TimeInterval{
				start:   2,
				end:     0,
				reverse: true,
			},
			minutes: minutesInDay - 2,
		},
		{
			name: "two reversed with one collisions",
			h: &TimeInterval{
				start:   15,
				end:     13,
				reverse: true,
			},
			other: &TimeInterval{
				start:   14,
				end:     12,
				reverse: true,
			},
			want: &TimeInterval{
				start:   15,
				end:     12,
				reverse: true,
			},
			minutes: minutesInDay - 3,
		},
		{
			name: "two reversed with no collisions",
			h: &TimeInterval{
				start:   15,
				end:     13,
				reverse: true,
			},
			other: &TimeInterval{
				start:   14,
				end:     12,
				reverse: true,
			},
			want: &TimeInterval{
				start:   15,
				end:     12,
				reverse: true,
			},
			minutes: minutesInDay - 3,
		},
		{
			name: "one reversed",
			h: &TimeInterval{
				start:   15,
				end:     13,
				reverse: true,
			},
			other: &TimeInterval{
				start: 10,
				end:   13,
			},
			want: &TimeInterval{
				start:   10,
				end:     13,
				reverse: false,
			},
			minutes: 4,
		},
		{
			name: "one reversed #2",
			other: &TimeInterval{
				start:   15,
				end:     13,
				reverse: true,
			},
			h: &TimeInterval{
				start: 10,
				end:   13,
			},
			want: &TimeInterval{
				start:   10,
				end:     13,
				reverse: false,
			},
			minutes: 4,
		},
		{
			name: "one reversed and has two collisions",
			other: &TimeInterval{
				start:   15,
				end:     13,
				reverse: true,
			},
			h: &TimeInterval{
				start: 10,
				end:   16,
			},
			want:    &TimeInterval{},
			minutes: 0,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, minutes := tc.h.Common(tc.other)
			assert.NotNil(t, got)
			assert.Equal(t, tc.minutes, minutes)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTimeInterval_Common_Negative_NilReference(t *testing.T) {
	assert.Panics(t, func() {
		(*TimeInterval)(nil).Common(nil)
	})
	assert.Panics(t, func() {
		(*TimeInterval)(nil).Common(new(TimeInterval))
	})
	assert.Panics(t, func() {
		(&TimeInterval{}).Common(nil)
	})
}

func TestTimeInterval_Duration(t *testing.T) {
	tt := []struct {
		name     string
		interval *TimeInterval
		want     Minute
	}{
		{"nil interval", nil, 0},
		{"reversed interval", &TimeInterval{14, 13, true}, minutesInDay - 1},
		{"non reversed", &TimeInterval{13, 14, false}, 2},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.interval.Duration())
		})
	}
}
