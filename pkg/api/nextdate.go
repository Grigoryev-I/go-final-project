package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	return date.Format(dateFormat) > now.Format(dateFormat)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения отсутствует")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {
	case "d":
		date, err = nextDay(date, now, parts)
	case "y":
		date, err = nextYear(date, now, parts)
	case "w":
		date, err = nextWeek(date, now, parts)
	case "m":
		date, err = nextMonth(date, now, parts)
	default:
		return "", fmt.Errorf("некорректный формат: %q", repeat)
	}

	if err != nil {
		return "", err
	}

	return date.Format(dateFormat), nil
}

// Повтор с интервалом
func nextDay(date, now time.Time, parts []string) (time.Time, error) {
	if len(parts) != 2 {
		return date, fmt.Errorf("некорректный формат для правила 'd'")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return date, fmt.Errorf("некорректный интервал: %q", parts[1])
	}

	if interval < 1 || interval > 400 {
		return date, fmt.Errorf("интервал дней от 1 до 400")
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}

	return date, nil
}

// Повтор ежегодный
func nextYear(date, now time.Time, parts []string) (time.Time, error) {
	if len(parts) != 1 {
		return date, fmt.Errorf("некорректный формат для правила 'y'")
	}

	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date, nil
}

// Повтор еженедельный
func nextWeek(date, now time.Time, parts []string) (time.Time, error) {
	if len(parts) != 2 {
		return date, fmt.Errorf("некорректный формат для правила 'w'")
	}

	var days [7]bool
	for _, s := range strings.Split(parts[1], ",") {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 7 {
			return date, fmt.Errorf("некорректный день: %q", s)
		}
		days[n-1] = true
	}

	for {
		date = date.AddDate(0, 0, 1)
		wd := int(date.Weekday())
		if wd == 0 {
			wd = 7
		}
		if afterNow(date, now) && days[wd-1] {
			break
		}
	}

	return date, nil
}

// Повтор ежемесячный
func nextMonth(date, now time.Time, parts []string) (time.Time, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return date, fmt.Errorf("некорректный формат для правила 'm'")
	}

	// дни месяца: 1-31, либо -1 (последний день), либо -2 (предпоследний)
	var days []int
	for _, s := range strings.Split(parts[1], ",") {
		n, err := strconv.Atoi(s)
		if err != nil || (n < 1 || n > 31) && n != -1 && n != -2 {
			return date, fmt.Errorf("некорректный день месяца: %q", s)
		}
		days = append(days, n)
	}

	// месяцы если не заданы, то подходят любые
	var months [13]bool
	allMonths := true
	if len(parts) == 3 {
		allMonths = false
		for _, s := range strings.Split(parts[2], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 12 {
				return date, fmt.Errorf("некорректный месяц: %q", s)
			}
			months[n-1] = true
		}
	}

	// перебор по одному дню, пока не найдём подходящий
	for i := 0; i < 366*10; i++ {
		date = date.AddDate(0, 0, 1)

		if !afterNow(date, now) {
			continue
		}
		if !allMonths && !months[int(date.Month())-1] {
			continue
		}
		if matchDay(date.Day(), daysInMonth(date), days) {
			return date, nil
		}
	}

	return date, fmt.Errorf("не удалось подобрать дату для правила 'm'")
}

// Проверяем, подходит ли день day под список дней days
func matchDay(day, monthLen int, days []int) bool {
	for _, d := range days {
		if d < 0 {
			d = monthLen + d + 1
		}
		if d == day {
			return true
		}
	}
	return false
}

// Возвращаем количество дней в месяце указанной даты
func daysInMonth(date time.Time) int {
	return time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now := time.Now()
	if nowStr != "" {
		parsed, err := time.Parse(dateFormat, nowStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now = parsed
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}
