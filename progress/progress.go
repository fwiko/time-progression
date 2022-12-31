package progress

import (
	"errors"
	"time"
)

type Values struct {
	Year   float64 `json:"year"`
	Month  float64 `json:"month"`
	Week   float64 `json:"week"`
	Day    float64 `json:"day"`
	Hour   float64 `json:"hour"`
	Minute float64 `json:"minute"`
}

type Response struct {
	Timezone  string `json:"timezone"`
	Timestamp string `json:"timestamp"`
	Result    Values `json:"result"`
}

const (
	SecondsInDay  = 86400
	SecondsInHour = 3600
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
	t := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	return float64(t.AddDate(0, 1, 0).Day())
}

func realWeekday(now time.Time) int {
	if now.Weekday() == time.Sunday {
		return 7
	} else {
		return int(now.Weekday())
	}
}

func seconds(t time.Time) Values {
	minuteSeconds := float64(t.Second())
	hourSeconds := float64(t.Minute()*60) + minuteSeconds
	daySeconds := float64(t.Hour()*SecondsInHour) + hourSeconds

	return Values{
		Year:   float64((t.YearDay()-1)*SecondsInDay) + daySeconds,
		Month:  float64((t.Day()-1)*SecondsInDay) + daySeconds,
		Week:   float64((realWeekday(t)-1)*SecondsInDay) + daySeconds,
		Day:    daySeconds,
		Hour:   hourSeconds,
		Minute: minuteSeconds,
	}
}

func percentages(t time.Time) Values {
	v := seconds(t)

	var daysInYear float64
	if isLeapYear(t.Year()) {
		daysInYear = 366
	} else {
		daysInYear = 365
	}

	daysInMonth := daysInMonth(t.Year(), t.Month())

	v.Year = v.Year / (daysInYear * float64(SecondsInDay)) * 100
	v.Month = v.Month / (daysInMonth * float64(SecondsInDay)) * 100
	v.Week = v.Week / (168 * float64(SecondsInHour)) * 100
	v.Day = v.Day / float64(SecondsInDay) * 100
	v.Hour = v.Hour / float64(SecondsInHour) * 100
	v.Minute = v.Minute / 60 * 100

	return v
}

func Query(format string, timezone string) (Response, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return Response{}, errors.New("invalid timezone")
	}

	now := time.Now().In(location)

	var result Values

	if format == "second" {
		result = seconds(now)
	} else if format == "percent" {
		result = percentages(now)
	} else {
		return Response{}, errors.New("invalid format")
	}

	return Response{
		Timezone:  timezone,
		Timestamp: now.Format(time.RFC1123),
		Result:    result,
	}, nil
}
