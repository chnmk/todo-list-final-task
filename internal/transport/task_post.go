package transport

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
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
	db, err := sql.Open("sqlite", DatabaseDir)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	id, err := res.LastInsertId()
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

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
