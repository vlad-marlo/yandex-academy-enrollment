package datetime

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const layout = "2006-01-02T15:04:05.000Z07:00"

var ErrNilReference = errors.New("nil reference in pointer receiver")

type Time time.Time

func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, ErrNilReference
	}
	s := time.Time(*t).Format(layout)
	return json.Marshal(s)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	if t == nil {
		return ErrNilReference
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("json: unmarshal: %w", err)
	}
	tm, err := time.Parse(layout, s)
	if err != nil {
		return fmt.Errorf("time: parse: %w", err)
	}
	*t = Time(tm)
	return nil
}

func (t *Time) Time() time.Time {
	if t == nil {
		return time.Time{}
	}
	return *(*time.Time)(t)
}
