package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func taskGET(w http.ResponseWriter, r *http.Request) {
	var task Task

	id := r.FormValue("id")
	if id == "" {
		returnError(w, "не указан идентификатор", 400)
		return
	}

	// Подключение к базе
	db, err := sql.Open("sqlite", DatabaseDir)
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	defer db.Close()

	// Выполнение запроса
	rows, err := db.Query("SELECT * FROM scheduler WHERE id = :id LIMIT 1", sql.Named("id", id))
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	// Чтение полученных данных
	for rows.Next() {
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			returnError(w, err.Error(), 500)
			return
		}
	}

	// Запись ответа
	resp, err := json.Marshal(task)
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}