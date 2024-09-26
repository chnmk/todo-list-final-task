package transport

import (
	"encoding/json"
	"net/http"

	"github.com/chnmk/todo-list-final-task/internal/database"
)

func taskDELETE(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		ReturnError(w, "не указан идентификатор", 400)
		return
	}

	// Удаление записи без правил повторения
	err := database.DeleteTaskById(DatabaseFile, id)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	// Если всё правильно, возвращает пустой ответ
	resp, err := json.Marshal(struct{}{})
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
