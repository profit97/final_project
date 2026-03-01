package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"il.karabach/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, http.StatusBadRequest, "не правильный JSON")
		return
	}

	if task.Title == "" {
		writeJSONError(w, http.StatusBadRequest, "пустое поле названия")
		return
	}

	now := time.Now()

	// Если дата не указана, ставим сегодняшнюю
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	// Проверка формата даты
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "не правильный формат даты")
		return
	}

	// Если задано правило повторения – проверяем его корректность и получаем следующую дату
	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Приводим даты к UTC и началу суток для сравнения
	todayUTC := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	taskDateUTC := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	// Если дата задачи меньше сегодняшней – корректируем
	if taskDateUTC.Before(todayUTC) {
		if task.Repeat == "" {
			// Без правила – ставим сегодня
			task.Date = now.Format(DateFormat)
		} else {
			// С правилом – используем вычисленную следующую дату
			task.Date = next
		}
	}

	// Добавляем в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "ошибка БД")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}