package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDate_Positive(t *testing.T) {
	date, err := ParseDate("2022-12-22")
	assert.NoError(t, err)
	if assert.NotNil(t, date) {
		assert.Equal(t, &Date{year: 2022, month: December, day: 22}, date)
	}
}

func TestParseDate_Negative(t *testing.T) {
	tt := []struct {
		name  string
		input string
	}{
		{"bad count of \"-\"", "2022-22"},
		{"year not parsable", "202d-12-31"},
		{"bad len of year", "202-12-31"},
		{"bad len of month", "2022-2-31"},
		{"bad len of day", "2022-12-1"},
		{"month not parsable", "2022-ad-32"},
		{"day not parsable", "2022-22-3s"},
		{"not valid", "2022-22-33"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			date, err := ParseDate(tc.input)
			assert.Nil(t, date)
			if assert.Error(t, err) {
				assert.ErrorIs(t, err, ErrBadDateFormat)
			}
		})
	}
}

func TestParseDate_valid(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.False(t, (*Date)(nil).valid())
	})
	tt := []struct {
		name  string
		year  int
		month int
		day   int
		want  assert.BoolAssertionFunc
	}{
		{"positive", 2002, December, 31, assert.True},
		{"bad day", 2002, December, 32, assert.False},
		{"bad day", 2002, April, 31, assert.False},
		{"bad month", 2002, -123, 32, assert.False},
		{"bad month", 2002, -123, -1, assert.False},
		{"non leap feb", 2100, February, 29, assert.False},
		{"non leap feb", 2000, February, 29, assert.True},
		{"leap feb", 2004, February, 29, assert.True},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			date := &Date{year: tc.year, month: tc.month, day: tc.day}
			tc.want(t, date.valid())
		})
	}
}

func TestIsYearLeap(t *testing.T) {
	for i := 0; i < 2022; i++ {
		if i%4 == 0 {
			if i%100 == 0 {
				if i%400 == 0 {
					assert.True(t, isYearLeap(i))
				} else {
					assert.False(t, isYearLeap(i))
				}
			} else {
				assert.True(t, isYearLeap(i))
			}
		} else {
			assert.False(t, isYearLeap(i))
		}
	}
}

func TestDate_String(t *testing.T) {
	tt := []struct {
		name string
		date Date
		want string
	}{
		{name: "empty date", want: "0000-00-00"},
		{"positive #1", Date{2022, 11, 30}, "2022-11-30"},
		{"positive #1", Date{2022, 11, 3}, "2022-11-03"},
		{"positive #1", Date{2022, 1, 30}, "2022-01-30"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.date.String())
			if !tc.date.valid() {
				return
			}
			date, err := ParseDate(tc.date.String())
			assert.NoError(t, err)
			if assert.NotNil(t, date) {
				assert.Equal(t, tc.date, *date)
			}
		})
	}
}

func TestDate_End_Start(t *testing.T) {
	todayEnd := Today().End()
	todayStart := Today().Start()
	nilStart := (*Date)(nil).Start()
	nilEnd := (*Date)(nil).End()
	assert.Equal(t, todayEnd, nilEnd)
	assert.Equal(t, todayStart, nilStart)
}

func TestDate_Less(t *testing.T) {
	tt := []struct {
		name  string
		date  *Date
		other *Date
		want  assert.BoolAssertionFunc
	}{
		{"all nil", nil, nil, assert.False},
		{"only other nil", new(Date), nil, assert.False},
		{"only date nil", nil, new(Date), assert.True},
		{"equal dates", &Date{2022, 12, 22}, &Date{2022, 12, 22}, assert.False},
		{"non equal years #1", &Date{2022, 12, 22}, &Date{2021, 12, 22}, assert.False},
		{"non equal years #2", &Date{2022, 12, 22}, &Date{2023, 12, 22}, assert.True},
		{"non equal months", &Date{2022, 11, 22}, &Date{2022, 10, 22}, assert.False},
		{"non equal months", &Date{2022, 11, 22}, &Date{2022, 12, 22}, assert.True},
		{"non equal days", &Date{2022, 11, 21}, &Date{2022, 11, 22}, assert.True},
		{"non equal days", &Date{2022, 11, 22}, &Date{2022, 11, 21}, assert.False},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.want(t, tc.date.Less(tc.other))
		})
	}
}
