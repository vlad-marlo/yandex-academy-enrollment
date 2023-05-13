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

// Minute is time implementation with underlying hour and minute
type Minute int

// ParseTime parses time in HH:MM format from string.
//
// Minute must be presented in HH:MM format.
//
// Only error that can be returned from this function is ErrBadWorkingHours of child
// error of itself.
func ParseTime(raw string) (w Minute, err error) {
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
	if hour > maxHourValue || hour < 0 {
		return w, fmt.Errorf("%w: hour must be less then 24", ErrBadWorkingHours)
	}

	if minute, err = strconv.Atoi(t[1]); err != nil {
		return w, fmt.Errorf("%w: err while parsing minute: %v", ErrBadWorkingHours, err)
	}
	if minute > maxMinuteValue || minute < 0 {
		return w, ErrBadWorkingHours
	}

	w = w.Add(hour*60 + minute)
	return w, nil
}

// In return is time in interval or not.
func (t Minute) In(h *TimeInterval) bool {
	if h.reverse {
		return t >= h.start || t <= h.end
	}
	return t >= h.start && t <= h.end
}

// String returns string representation of time.
//
// Minute will be in format HH:MM and must be parsable back to Minute object from string.
func (t Minute) String() string {
	var hour, minute string
	h, m := t.Hour(), t.Minute()
	if h < 10 {
		hour = fmt.Sprintf("0%d", h)
	} else {
		hour = fmt.Sprint(h)
	}
	if m < 10 {
		minute = fmt.Sprintf("0%d", m)
	} else {
		minute = fmt.Sprint(m)
	}
	return fmt.Sprintf("%s:%s", hour, minute)
}

// Hour is hour accessor.
func (t Minute) Hour() int {
	return int(t) / 60
}

// Minute is minute accessor.
func (t Minute) Minute() int {
	return int(t) % 60
}

func (t Minute) Add(minutes int) Minute {
	return Minute((int(t) + minutes) % (24 * 60))
}
