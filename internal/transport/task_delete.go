package transport

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func taskDELETE(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		ReturnError(w, "не указан идентификатор", 400)
		return
	}

	// Получение нужной записи
	db, err := sql.Open("sqlite", DatabaseDir)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	defer db.Close()

	// Удаление записи без правил повторения
	del, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id),
	)
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	}

	rows, err := del.RowsAffected()
	if err != nil {
		ReturnError(w, err.Error(), 500)
		return
	} else if rows == 0 {
		ReturnError(w, "задача не найдена", 500)
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
