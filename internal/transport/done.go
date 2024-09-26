package transport

import (
	"encoding/json"
	"net/http"

	"github.com/chnmk/todo-list-final-task/internal/database"
)

// Удаляет или переносит на следующую дату задачу с полученном в запросе id.
func TaskDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ReturnError(w, "неожиданный метод запроса", 500)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		ReturnError(w, "не указан идентификатор", 400)
		return
	}

	err := database.CompleteTaskById(DatabaseFile, id)
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
