package transport

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/chnmk/todo-list-final-task/internal/database"
)

func taskPUT(w http.ResponseWriter, r *http.Request) {
	task, err := checkRequest(w, r)
	if err != nil {
		return
	}
	if task.Id == "" {
		ReturnError(w, "не указан id задачи", 400)
		return
	}

	// Если всё правильно, редактирует запись
	err = database.UpdateTask(DatabaseFile, task)
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

	log.Println("Success!")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
