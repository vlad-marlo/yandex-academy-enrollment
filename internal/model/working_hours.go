package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrBadWorkingHours = errors.New("working hours must be in HH:MM-HH:MM format")
	//workingHoursRegexp = regexp.MustCompile("^\\d+:\\d+-\\d+:\\d+$")
)

const (
	maxHourValue   = 23
	maxMinuteValue = 59
)

type (
	// Time is time implementation with underlying hour and minute
	Time struct {
		hour   uint8
		minute uint8
	}
	// TimeInterval represents time interval.
	TimeInterval struct {
		start   Time
		end     Time
		reverse bool
	}
)

// ParseWorkTime parses time from raw string.
//
// Time must be presented in HH:MM format.
//
// Only error that can be returned from this function is ErrBadWorkingHours of child
// error of itself.
func ParseWorkTime(raw string) (w Time, err error) {
	time := strings.Split(raw, ":")
	if len(time) != 2 {
		return Time{}, ErrBadWorkingHours
	}
	var hour, minute int
	if hour, err = strconv.Atoi(time[0]); err != nil {
		return w, fmt.Errorf("%w: err while parsing hour: %v", ErrBadWorkingHours, err)
	}
	if hour > maxHourValue {
		return w, ErrBadWorkingHours
	}

	if minute, err = strconv.Atoi(time[1]); err != nil {
		return w, fmt.Errorf("%w: err while parsing minute: %v", ErrBadWorkingHours, err)
	}
	if minute > maxMinuteValue {
		return w, ErrBadWorkingHours
	}

	w.hour = uint8(hour)
	w.minute = uint8(minute)
	return w, nil
}

// ParseWorkingHours parses working hours from string.
//
// Time must be provided in HH:MM-HH:MM format.
// If string will not match this pattern then will be returned ErrBadWorkingHours.
//
// Only error that can be returned from this function is ErrBadWorkingHours of child
// error of itself.
func ParseWorkingHours(raw string) (h *TimeInterval, err error) {
	hours := strings.Split(raw, "-")
	if len(hours) != 2 {
		return nil, ErrBadWorkingHours
	}

	h = new(TimeInterval)
	if h.start, err = ParseWorkTime(hours[0]); err != nil {
		return nil, err
	}
	if h.end, err = ParseWorkTime(hours[1]); err != nil {
		return nil, err
	}
	h.reverse = h.end.Less(h.start)

	return h, nil
}

// In return is time in interval or not.
func (h *TimeInterval) In(time Time) bool {
	if h.reverse {
		return (h.start.Less(time) && !(time.Less(h.end))) || time.Less(h.end)
	}
	return time.Less(h.end) && h.start.Less(time)
}

// Less is helper comparator function.
func (h Time) Less(other Time) bool {
	if h.hour == other.hour {
		return h.minute < other.minute
	}
	return h.hour < other.hour
}
