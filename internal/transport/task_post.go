package transport

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/chnmk/todo-list-final-task/internal/database"
)

func taskPOST(w http.ResponseWriter, r *http.Request) {
	var responseValid ResponseValid
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Считывает данные
	task, err := checkRequest(w, r)
	if err != nil {
		return
	}

	// Если всё правильно, добавляет новую запись
	id, err := database.AddTask(DatabaseFile, task)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	// Если всё правильно, возвращает ответ
	responseValid.Id = strconv.FormatInt(id, 10)
	resp, err := json.Marshal(responseValid)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	log.Println("Success!")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
