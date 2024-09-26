package transport

import (
	"encoding/json"
	"net/http"

	"github.com/chnmk/todo-list-final-task/internal/database"
	"github.com/chnmk/todo-list-final-task/internal/services"
)

func taskGET(w http.ResponseWriter, r *http.Request) {
	var task services.Task

	id := r.FormValue("id")
	if id == "" {
		ReturnError(w, "не указан идентификатор", 400)
		return
	}

	task, err := database.GetTaskById(DatabaseFile, id)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	// Запись ответа
	resp, err := json.Marshal(task)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
