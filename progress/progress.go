package progress

import (
	"errors"
	"time"
)

type values struct {
	Year   float64 `json:"year"`
	Month  float64 `json:"month"`
	Week   float64 `json:"week"`
	Day    float64 `json:"day"`
	Hour   float64 `json:"hour"`
	Minute float64 `json:"minute"`
}

type query struct {
	Timezone  string `json:"timezone"`
	Timestamp string `json:"timestamp"`
	Result    values `json:"result"`
}

var (
	SECONDS_IN_DAY  = 86400
	SECONDS_IN_HOUR = 3600
	HOURS_IN_WEEK   = 168
)

func isLeapYear(year int) bool {
	if year%4 == 0 && year%100 != 0 {
		return true
	} else if year%400 == 0 {
		return true
	}
	return false
}

func daysInMonth(year int, month time.Month) float64 {
	t := time.Date(year, time.Month(month), 0, 0, 0, 0, 0, time.UTC)
	return float64(t.AddDate(0, 1, 0).Day())
}

func realWeekday(now time.Time) int {
	if now.Weekday() == time.Sunday {
		return 7
	} else {
		return int(now.Weekday())
	}
}

func seconds(t time.Time) values {
	minuteSeconds := float64(t.Second())
	hourSeconds := float64(t.Minute()*60) + minuteSeconds
	daySeconds := float64(t.Hour()*SECONDS_IN_HOUR) + hourSeconds

	return values{
		Year:   float64(t.YearDay()*SECONDS_IN_DAY) + daySeconds,
		Month:  float64(t.Day()*SECONDS_IN_DAY) + daySeconds,
		Week:   float64((realWeekday(t)-1)*SECONDS_IN_DAY) + daySeconds,
		Day:    daySeconds,
		Hour:   hourSeconds,
		Minute: minuteSeconds,
	}
}

func percentages(t time.Time) values {
	v := seconds(t)

	var daysInYear float64
	if isLeapYear(t.Year()) {
		daysInYear = 366
	} else {
		daysInYear = 365
	}

	daysInMonth := daysInMonth(t.Year(), t.Month())

	v.Year = v.Year / (daysInYear * float64(SECONDS_IN_DAY)) * 100
	v.Month = v.Month / (daysInMonth * float64(SECONDS_IN_DAY)) * 100
	v.Week = v.Week / (168 * float64(SECONDS_IN_HOUR)) * 100
	v.Day = v.Day / float64(SECONDS_IN_DAY) * 100
	v.Hour = v.Hour / float64(SECONDS_IN_HOUR) * 100
	v.Minute = v.Minute / 60 * 100

	return v
}

func Query(format string, timezone string) (query, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return query{}, errors.New("Invalid timezone")
	}

	now := time.Now().In(location)

	var result values

	if format == "second" {
		result = seconds(now)
	} else if format == "percentage" {
		result = percentages(now)
	} else {
		return query{}, errors.New("Invalid format")
	}

	return query{
		Timezone:  timezone,
		Timestamp: now.Format(time.RFC1123),
		Result:    result,
	}, nil
}
