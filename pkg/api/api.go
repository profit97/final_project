package api

import (
	"encoding/json"
	"net/http"
	"time"

	"il.karabach/pkg/db"
)

const DateFormat = "20060102"

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// getTaskHandler handles GET /api/task?id=...
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "пропущен id параметр")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err.Error() == "задача не найдена" {
			writeJSONError(w, http.StatusNotFound, "задача не найдена")
		} else {
			writeJSONError(w, http.StatusInternalServerError, "ошибка БД")
		}
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// updateTaskHandler handles PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, http.StatusBadRequest, "неправильный JSON")
		return
	}

	if task.ID == "" {
		writeJSONError(w, http.StatusBadRequest, "пропущено поле id")
		return
	}
	if task.Title == "" {
		writeJSONError(w, http.StatusBadRequest, "пустое поле заголовка")
		return
	}
	if task.Date == "" {
		writeJSONError(w, http.StatusBadRequest, "пустое поле даты")
		return
	}

	_, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "не правильный формат даты")
		return
	}

	if task.Repeat != "" {
		now := time.Now()
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	err = db.UpdateTask(&task)
	if err != nil {
		if err.Error() == "задача не найдена" {
			writeJSONError(w, http.StatusNotFound, "задача не найдена")
		} else {
			writeJSONError(w, http.StatusInternalServerError, "ошибка БД")
		}
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}

// deleteTaskHandler handles DELETE /api/task?id=...
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "пропущен параметр id")
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if err.Error() == "задача не найдена" {
			writeJSONError(w, http.StatusNotFound, "задача не найдена")
		} else {
			writeJSONError(w, http.StatusInternalServerError, "ошибка БД")
		}
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}

// taskHandler routes requests based on method
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Init регистрирует все обработчики
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}