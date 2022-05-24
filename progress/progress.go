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

func Query(timezone string) (query, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return query{}, errors.New("Invalid timezone")
	}

	now := time.Now().In(location)

	minuteSeconds := float64(now.Second())
	hourSeconds := float64(now.Minute()*60) + minuteSeconds
	daySeconds := float64(now.Hour()*SECONDS_IN_HOUR) + hourSeconds
	weekSeconds := float64((realWeekday(now)-1)*SECONDS_IN_DAY) + daySeconds
	monthSeconds := float64(now.Day()*SECONDS_IN_DAY) + daySeconds
	yearSeconds := float64(now.YearDay()*SECONDS_IN_DAY) + daySeconds

	var daysInYear float64
	if isLeapYear(now.Year()) {
		daysInYear = 366
	} else {
		daysInYear = 365
	}

	daysInMonth := daysInMonth(now.Year(), now.Month())

	return query{
		Timezone:  timezone,
		Timestamp: now.Format(time.RFC1123),
		Result: values{
			Year:   yearSeconds / (daysInYear * float64(SECONDS_IN_DAY)) * 100,
			Month:  monthSeconds / (daysInMonth * float64(SECONDS_IN_DAY)) * 100,
			Week:   weekSeconds / (168 * float64(SECONDS_IN_HOUR)) * 100,
			Day:    daySeconds / float64(SECONDS_IN_DAY) * 100,
			Hour:   hourSeconds / float64(SECONDS_IN_HOUR) * 100,
			Minute: minuteSeconds / 60 * 100,
		},
	}, nil
}
