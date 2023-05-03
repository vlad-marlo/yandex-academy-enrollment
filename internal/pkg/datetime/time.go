package datetime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrBadWorkingHours = errors.New("working hours must be in HH:MM-HH:MM format")
	ErrBadDateFormat   = errors.New("date format must be in YYYY-MM-DD format")
)

const (
	maxHourValue   = 23
	maxMinuteValue = 59
)

// Time is time implementation with underlying hour and minute
type Time struct {
	hour   uint8
	minute uint8
}

// ParseTime parses time in HH:MM format from string.
//
// Time must be presented in HH:MM format.
//
// Only error that can be returned from this function is ErrBadWorkingHours of child
// error of itself.
func ParseTime(raw string) (w Time, err error) {
	t := strings.Split(raw, ":")
	if len(t) != 2 {
		return w, ErrBadWorkingHours
	}
	var hour, minute int
	if hour, err = strconv.Atoi(t[0]); err != nil {
		return w, fmt.Errorf("%w: err while parsing hour: %v", ErrBadWorkingHours, err)
	}
	if len(t[0]) != 2 || len(t[1]) != 2 {
		return w, ErrBadWorkingHours
	}
	if hour > maxHourValue {
		return w, fmt.Errorf("%w: hour must be less then 24", ErrBadWorkingHours)
	}

	if minute, err = strconv.Atoi(t[1]); err != nil {
		return w, fmt.Errorf("%w: err while parsing minute: %v", ErrBadWorkingHours, err)
	}
	if minute > maxMinuteValue {
		return w, ErrBadWorkingHours
	}

	w.hour = uint8(hour)
	w.minute = uint8(minute)
	return w, nil
}

// In return is time in interval or not.
func (t Time) In(h *TimeInterval) bool {
	if h.reverse {
		return (h.start.Less(t) && !(t.Less(h.end))) || t.Less(h.end)
	}
	return t.Less(h.end) && h.start.Less(t)
}

// String returns string representation of time.
//
// Time will be in format HH:MM and must be parsable back to Time object from string.
func (t Time) String() string {
	var hour, minute string
	if t.hour < 10 {
		hour = fmt.Sprintf("0%d", t.hour)
	} else {
		hour = fmt.Sprint(t.hour)
	}
	if t.minute < 10 {
		minute = fmt.Sprintf("0%d", t.minute)
	} else {
		minute = fmt.Sprint(t.minute)
	}
	return fmt.Sprintf("%s:%s", hour, minute)
}

// Hour is hour accessor.
func (t Time) Hour() uint8 {
	return t.hour
}

// Minute is minute accessor.
func (t Time) Minute() uint8 {
	return t.minute
}

// Less is helper comparator function.
func (t Time) Less(other Time) bool {
	if t.hour == other.hour {
		return t.minute < other.minute
	}
	return t.hour < other.hour
}

func (t Time) Add(minutes uint8) Time {
	res := Time{
		minute: (t.minute + minutes) % 60,
		hour:   (t.hour + minutes/60) % 24,
	}
	if t.minute+(minutes%60) >= 60 {
		res.hour = (res.hour + 1) % 24
	}
	return res
}
