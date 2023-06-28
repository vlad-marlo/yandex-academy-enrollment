package datetime

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
	"time"
)

// TimeInterval represents time interval.
type TimeInterval struct {
	start   Minute
	end     Minute
	reverse bool
}

const (
	minutesInDay = Minute(24 * 60)
)

func min[T constraints.Integer | constraints.Float](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T constraints.Integer | constraints.Float](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// TimeIntervalAlias is alias struct for TimeInterval.
//
// It must be used only and only then you need to create TimeInterval with some data.
type TimeIntervalAlias struct {
	Start   int32
	End     int32
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
	h.reverse = h.end < h.start

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

func (h *TimeInterval) calculateCommonForReversed(other *TimeInterval) (*TimeInterval, Minute) {
	if !h.reverse {
		return other.calculateCommonForReversed(h)
	}
	if h.start <= other.end && h.end >= other.start {
		return new(TimeInterval), 0
	}

	start := min(h.end, other.start)
	end := min(h.start, other.end)

	res := &TimeInterval{start: start, end: end, reverse: start > end}
	return res, res.Duration()
}

func (h *TimeInterval) Common(other *TimeInterval) (*TimeInterval, Minute) {
	if h == nil || other == nil {
		panic("unexpectedly got nil pointer reference")
	}
	if h.reverse != other.reverse {
		return h.calculateCommonForReversed(other)
	}

	start := max(h.start, other.start)
	end := min(h.end, other.end)

	if (h.start > other.end || h.end < other.start) && !h.reverse {
		return new(TimeInterval), 0
	}

	res := &TimeInterval{start: start, end: end, reverse: start > end}
	return res, res.Duration()
}

func (h *TimeInterval) Duration() Minute {
	if h == nil {
		return 0
	}

	if h.reverse {
		return minutesInDay - h.start + h.end
	}
	return h.end - h.start + 1
}

// Start return start of time interval.
//
// If pointer receiver is nil object then will be returned zero time.
func (h *TimeInterval) Start() (t Minute) {
	if h == nil {
		return 0
	}
	return h.start
}

// End return end of time interval.
//
// If pointer receiver is nil object then will be returned zero time.
func (h *TimeInterval) End() (t Minute) {
	if h == nil {
		return 0
	}
	return h.end
}
