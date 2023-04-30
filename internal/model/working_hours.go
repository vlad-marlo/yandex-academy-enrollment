package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrBadWorkingHours = errors.New("working hours must be in HH:MM-HH:MM format")
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

// ParseWorkTime parses time in HH:MM format from string.
//
// Time must be presented in HH:MM format.
//
// Only error that can be returned from this function is ErrBadWorkingHours of child
// error of itself.
func ParseWorkTime(raw string) (w Time, err error) {
	t := strings.Split(raw, ":")
	if len(t) != 2 {
		return w, ErrBadWorkingHours
	}
	var hour, minute int
	if hour, err = strconv.Atoi(t[0]); err != nil {
		return w, fmt.Errorf("%w: err while parsing hour: %v", ErrBadWorkingHours, err)
	}
	if hour > maxHourValue {
		return w, ErrBadWorkingHours
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

// ParseWorkingHours parses working hours from string.
//
// Time must be provided in HH:MM-HH:MM format.
// If string will not match this pattern then will be returned ErrBadWorkingHours.
//
// Only error that can be returned from this function is ErrBadWorkingHours of child
// error of itself.
func ParseWorkingHours(raw string) (h *TimeInterval, err error) {
	raw = strings.Trim(raw, "\"")
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

// String returns string representation of time interval.
//
// String representation of time interval always parsable.
//
// Time interval is returned in HH:MM-HH:MM format.
func (h *TimeInterval) String() string {
	return fmt.Sprintf("%s-%s", h.start, h.end)
}

// TimeIn return is time in interval or not.
func (h *TimeInterval) TimeIn(t time.Time) bool {
	return Time{
		hour:   uint8(t.Hour()),
		minute: uint8(t.Minute()),
	}.In(h)
}

// MarshalJSON makes available to represent time interval in JSON format which is parsable back to TimeInterval obj.
func (h *TimeInterval) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON makes available to parse TimeInterval from JSON representation format.
func (h *TimeInterval) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	interval, err := ParseWorkingHours(raw)
	if err != nil {
		return err
	}
	*h = *interval
	return nil
}

// Start return start of time interval.
//
// If pointer receiver is nil object then will be returned zero time.
func (h *TimeInterval) Start() (t Time) {
	if h == nil {
		return
	}
	return h.start
}

// End return end of time interval.
//
// If pointer receiver is nil object then will be returned zero time.
func (h *TimeInterval) End() (t Time) {
	if h == nil {
		return
	}
	return h.end
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
	return fmt.Sprintf("%d:%d", t.hour, t.minute)
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
