package datetime

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// TimeInterval represents time interval.
type TimeInterval struct {
	start   Minute
	end     Minute
	reverse bool
}

// TimeIntervalAlias is alias struct for TimeInterval.
//
// It must be used only and only then you need to create TimeInterval with some data.
type TimeIntervalAlias struct {
	Start   int
	End     int
	Reverse bool
}

// ParseTimeInterval parses working hours from string.
//
// Minute must be provided in HH:MM-HH:MM format.
// If string will not match this pattern then will be returned ErrBadWorkingHours.
//
// Only error that can be returned from this function is ErrBadWorkingHours of child
// error of itself.
func ParseTimeInterval(raw string) (h *TimeInterval, err error) {
	raw = strings.Trim(raw, "\"")
	hours := strings.Split(raw, "-")
	if len(hours) != 2 {
		return nil, ErrBadWorkingHours
	}

	h = new(TimeInterval)
	if h.start, err = ParseTime(hours[0]); err != nil {
		return nil, err
	}
	if h.end, err = ParseTime(hours[1]); err != nil {
		return nil, err
	}
	h.reverse = h.end.Less(h.start)

	return h, nil
}

func (t TimeIntervalAlias) TimeInterval() *TimeInterval {
	return &TimeInterval{
		start:   Minute(t.Start),
		end:     Minute(t.End),
		reverse: t.Reverse,
	}
}

// String returns string representation of time interval.
//
// String representation of time interval always parsable.
//
// Minute interval is returned in HH:MM-HH:MM format.
func (h *TimeInterval) String() string {
	return fmt.Sprintf("%s-%s", h.start, h.end)
}

// TimeIn return is time in interval or not.
func (h *TimeInterval) TimeIn(t time.Time) bool {
	return Minute(60*t.Hour() + t.Minute()).In(h)
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
	interval, err := ParseTimeInterval(raw)
	if err != nil {
		return err
	}
	*h = *interval
	return nil
}

// Start return start of time interval.
//
// If pointer receiver is nil object then will be returned zero time.
func (h *TimeInterval) Start() (t Minute) {
	if h == nil {
		return
	}
	return h.start
}

// End return end of time interval.
//
// If pointer receiver is nil object then will be returned zero time.
func (h *TimeInterval) End() (t Minute) {
	if h == nil {
		return
	}
	return h.end
}
