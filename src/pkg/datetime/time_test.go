package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	TestTime1 = Minute(11*60 + 59)
	TestTime2 = Minute(12*60 + 1)
)

func TestTime_In(t *testing.T) {
	t.Run("non reversed", func(t *testing.T) {
		h := &TimeInterval{
			start:   TestTime1,
			end:     TestTime2,
			reverse: false,
		}
		assert.True(t, TestTime1.In(h))
		assert.True(t, TestTime2.In(h))
		assert.False(t, (TestTime1 - 1).In(h))
		assert.False(t, (TestTime2 + 1).In(h))
	})
	t.Run("reversed", func(t *testing.T) {
		h := &TimeInterval{
			start:   TestTime2,
			end:     TestTime1,
			reverse: true,
		}
		assert.True(t, TestTime1.In(h))
		assert.True(t, TestTime2.In(h))
		assert.False(t, (TestTime1 + 1).In(h))
		assert.False(t, (TestTime2 - 1).In(h))
	})
}

func TestParseTime(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		h, err := ParseTime("12:02")
		assert.NoError(t, err)
		assert.Equal(t, 12, h.Hour())
		assert.Equal(t, 2, h.Minute())
		parseFromH, err := ParseTime(h.String())
		assert.NoError(t, err)
		assert.Equal(t, h, parseFromH)
	})
	t.Run("bad string format", func(t *testing.T) {
		h, err := ParseTime("1:2:22")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Minute(0), h)
	})
	t.Run("hour is not parsable", func(t *testing.T) {
		h, err := ParseTime("d:12")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Minute(0), h)
	})
	t.Run("hour must be two-digit int (00, 01, 02 etc.)", func(t *testing.T) {
		h, err := ParseTime("1:12")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Minute(0), h)
	})
	t.Run("minute must be two-digit int (00, 01, 02 etc.)", func(t *testing.T) {
		h, err := ParseTime("12:2")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Minute(0), h)
	})
	t.Run("minute is not parsable", func(t *testing.T) {
		h, err := ParseTime("11:sdf")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Minute(0), h)
	})
	t.Run("hour must be less then 24", func(t *testing.T) {
		h, err := ParseTime("24:22")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Minute(0), h)
	})
	t.Run("minute must be less then 24", func(t *testing.T) {
		h, err := ParseTime("22:60")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Minute(0), h)
	})
}

func TestTime_Hour(t *testing.T) {
	assert.Equal(t, int(TestTime1)/60, TestTime1.Hour())
	assert.Equal(t, int(TestTime2)/60, TestTime2.Hour())
}

func TestTime_Minute(t *testing.T) {
	assert.Equal(t, int(TestTime1)%60, TestTime1.Minute())
	assert.Equal(t, int(TestTime2)%60, TestTime2.Minute())
}

func TestTime_Add(t *testing.T) {
	tt := []struct {
		name     string
		before   Minute
		expected Minute
		add      int
	}{
		{"default", Minute(12*60 + 22), Minute(13*60 + 40), 78},
		{"minutes overflows", Minute(12*60 + 22), Minute(13*60 + 00), 38},
		{"hours overflows", Minute(23*60 + 59), Minute(2*60 + 59), 180},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.before.Add(tc.add))
		})
	}
}

func TestTime_String(t *testing.T) {
	tt := []struct {
		name   string
		hour   int
		minute int
		want   string
	}{
		{"hour >= 10 and minute >= 10", 10, 10, "10:10"},
		{"hour < 10 and minute >= 10", 9, 10, "09:10"},
		{"hour >= 10 and minute < 10", 10, 9, "10:09"},
		{"hour < 10 and minute < 10", 9, 9, "09:09"},
		{"zero values", 0, 0, "00:00"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, Minute(tc.hour*60+tc.minute).String())
		})
	}
}

func TestParseTime_Negative(t *testing.T) {
	tt := []struct {
		name string
		raw  string
	}{
		{"bad count of \":\"", "123:123:123"},
		{"hour is not parsable", "12d:23"},
		{"bad len of hour", "2:22"},
		{"bad len of minute", "22:2"},
		{"hour is bigger then max val", "24:22"},
		{"hour is lower then min val", "-1:22"},
		{"minute is not parsable", "22:as"},
		{"minute is bigger then max val", "22:60"},
		{"minute is lower then min val", "22:-1"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseTime(tc.raw)
			if assert.Error(t, err) {
				assert.ErrorIs(t, err, ErrBadWorkingHours)
			}
		})
	}
}
