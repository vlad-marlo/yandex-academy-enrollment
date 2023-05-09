package datetime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	January = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
	maxDay = 31
)

type Date struct {
	year  int
	month int
	day   int
}

func Today() (d *Date) {
	d = new(Date)
	var month time.Month
	d.year, month, d.day = time.Now().In(time.UTC).Date()
	d.month = int(month)
	return
}

// ParseDate parses date from raw string.
//
// Date must be in format YEAR-MONTH-DAY
func ParseDate(raw string) (d *Date, err error) {
	d = new(Date)
	split := strings.Split(raw, "-")
	if len(split) != 3 {
		return nil, ErrBadDateFormat
	}
	if d.year, err = strconv.Atoi(split[0]); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadDateFormat, err)
	}
	if len(split[0]) != 4 || len(split[1]) != 2 || len(split[2]) != 2 {
		return nil, ErrBadDateFormat
	}
	if d.month, err = strconv.Atoi(split[1]); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadDateFormat, err)
	}
	if d.day, err = strconv.Atoi(split[2]); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBadDateFormat, err)
	}
	if !d.valid() {
		return nil, ErrBadDateFormat
	}

	return
}

func (date *Date) valid() bool {
	if date == nil {
		return false
	}
	if date.day < 1 {
		return false
	}
	switch date.month {
	case January, March, May, July, August, October, December:
		return date.day <= maxDay
	case February:
		return (date.day == 29 && isYearLeap(date.year)) || (date.day <= 28)
	case April, June, September, November:
		return date.day < maxDay
	default:
		return false
	}
}

func (date *Date) String() string {
	var (
		year  = fmt.Sprint(date.year)
		month = fmt.Sprint(date.month)
		day   = fmt.Sprint(date.day)
	)
	year = fmt.Sprintf("%s%s", strings.Repeat("0", 4-len(year)), year)
	month = fmt.Sprintf("%s%s", strings.Repeat("0", 2-len(month)), month)
	day = fmt.Sprintf("%s%s", strings.Repeat("0", 2-len(day)), day)
	return fmt.Sprintf("%s-%s-%s", year, month, day)
}

func isYearLeap(year int) bool {
	return !(year%4 != 0 || (year%100 == 0 && year%400 != 0))
}
