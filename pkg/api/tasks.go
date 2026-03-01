package api

import (
	"net/http"

	"il.karabach/pkg/db"
)

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "ошибка базы данных")
		return
	}

	writeJSON(w, http.StatusOK, tasksResponse{Tasks: tasks})
}