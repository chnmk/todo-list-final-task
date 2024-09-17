package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func taskPUT(w http.ResponseWriter, r *http.Request) {
	task, err := checkRequest(w, r)
	if err != nil {
		return
	}
	if task.Id == "" {
		returnError(w, "не указан id задачи", 400)
		return
	}

	// Если всё правильно, редактирует запись
	db, err := sql.Open("sqlite", DatabaseDir)
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	defer db.Close()

	temp, err := db.Exec("UPDATE scheduler SET (date, title, comment, repeat) = (:date, :title, :comment, :repeat) WHERE id = :id",
		sql.Named("id", task.Id),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	rows, err := temp.RowsAffected()
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	} else if rows == 0 {
		returnError(w, "задача не найдена", 500)
		return
	}

	// Если всё правильно, возвращает пустой ответ
	resp, err := json.Marshal(ResponseValid{})
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}