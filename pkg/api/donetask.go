package api

import (
	"net/http"
	"time"

	"il.karabach/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "пропущен параметр")
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

	if task.Repeat == "" {
		// Одноразовая задача - удаляем
		err = db.DeleteTask(id)
	} else {
		// Периодическая задача - вычисляем следующую дату и обновляем
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		err = db.UpdateTaskDate(next, id)
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "ошибка БД")
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}