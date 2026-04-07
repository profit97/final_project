package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустое значение повторений")
	}

	start, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("не правильная дата начала: %w", err)
	}

	// Приводим обе даты к UTC и началу суток
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)

	switch {
	case strings.HasPrefix(repeat, "d "):
		parts := strings.Split(repeat, " ")
		if len(parts) != 2 {
			return "", errors.New("не правильный формат повтора: нужно 'd <days>'")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("не правильное чисдл: %w", err)
		}
		if days < 1 || days > 400 {
			return "", errors.New("дневной интервал должен быть от 1 до 400")
		}
		date := start
		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case repeat == "y":
		date := start
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	default:
		return "", errors.New("не поддерживаемый повтор")
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	if dateStr == "" {
		http.Error(w, "не выбрана дата", http.StatusBadRequest)
		return
	}
	if repeatStr == "" {
		http.Error(w, "не выбран параметр повтора", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "не правильный формат, ожидалось YYYYMMDD", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}