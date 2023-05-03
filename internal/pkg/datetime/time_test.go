package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	TestTime1 = Time{
		hour:   11,
		minute: 59,
	}
	TestTime2 = Time{
		hour:   12,
		minute: 01,
	}
)

func TestParseTime(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		h, err := ParseTime("12:02")
		assert.NoError(t, err)
		assert.Equal(t, uint8(12), h.Hour())
		assert.Equal(t, uint8(2), h.Minute())
		parseFromH, err := ParseTime(h.String())
		assert.NoError(t, err)
		assert.Equal(t, h, parseFromH)
	})
	t.Run("bad string format", func(t *testing.T) {
		h, err := ParseTime("1:2:22")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Time{}, h)
	})
	t.Run("hour is not parsable", func(t *testing.T) {
		h, err := ParseTime("d:12")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Time{}, h)
	})
	t.Run("hour must be two-digit int (00, 01, 02 etc.)", func(t *testing.T) {
		h, err := ParseTime("1:12")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Time{}, h)
	})
	t.Run("minute must be two-digit int (00, 01, 02 etc.)", func(t *testing.T) {
		h, err := ParseTime("12:2")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Time{}, h)
	})
	t.Run("minute is not parsable", func(t *testing.T) {
		h, err := ParseTime("11:sdf")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Time{}, h)
	})
	t.Run("hour must be less then 24", func(t *testing.T) {
		h, err := ParseTime("24:22")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Time{}, h)
	})
	t.Run("minute must be less then 24", func(t *testing.T) {
		h, err := ParseTime("22:60")
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, ErrBadWorkingHours)
		}
		assert.Equal(t, Time{}, h)
	})
}

func TestTime_Less(t *testing.T) {
	t.Run("different hours", func(t *testing.T) {
		t1, t2 := Time{hour: 12, minute: 23}, Time{hour: 14, minute: 24}
		assert.True(t, t1.Less(t2))
		assert.False(t, t2.Less(t1))
	})
	t.Run("same hours", func(t *testing.T) {
		t1, t2 := Time{hour: 12, minute: 23}, Time{hour: 12, minute: 24}
		assert.True(t, t1.Less(t2))
		assert.False(t, t2.Less(t1))
	})
}

func TestTime_Hour(t *testing.T) {
	assert.Equal(t, TestTime1.hour, TestTime1.Hour())
	assert.Equal(t, TestTime2.hour, TestTime2.Hour())
}

func TestTime_Minute(t *testing.T) {
	assert.Equal(t, TestTime1.minute, TestTime1.Minute())
	assert.Equal(t, TestTime2.minute, TestTime2.Minute())
}

func TestTime_Add(t *testing.T) {
	tt := []struct {
		name     string
		before   Time
		expected Time
		add      uint8
	}{
		{"default", Time{12, 22}, Time{13, 40}, 78},
		{"minutes overflows", Time{12, 22}, Time{13, 00}, 38},
		{"hours overflows", Time{23, 59}, Time{2, 59}, 180},
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
		hour   uint8
		minute uint8
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
			assert.Equal(t, tc.want, Time{hour: tc.hour, minute: tc.minute}.String())
		})
	}
}
