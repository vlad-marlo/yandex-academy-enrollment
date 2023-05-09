package datetime

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const layout = "2006-01-02T15:04:05.000Z07:00"

var ErrNilReference = errors.New("nil reference in pointer receiver")

type TimeAlias time.Time

func (t *TimeAlias) MarshalJSON() ([]byte, error) {
	if t == nil {
		return nil, ErrNilReference
	}
	s := time.Time(*t).Format(layout)
	return json.Marshal(s)
}

func (t *TimeAlias) UnmarshalJSON(data []byte) error {
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
	*t = TimeAlias(tm)
	return nil
}

func (t *TimeAlias) Time() time.Time {
	if t == nil {
		return time.Time{}
	}
	return *(*time.Time)(t)
}
