package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/chnmk/todo-list-final-task/services"
)

func TaskDone(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		DeleteTask(w, r)
		return
	} else if r.Method != http.MethodPost {
		returnError(w, "неожиданный метод запроса", 500)
		return
	}

	var task Task

	id := r.FormValue("id")
	if id == "" {
		returnError(w, "не указан идентификатор", 400)
		return
	}

	// Получение нужной записи
	db, err := sql.Open("sqlite", DatabaseDir)
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	defer db.Close()

	row := db.QueryRow("SELECT date, repeat FROM scheduler WHERE id = :id LIMIT 1", sql.Named("id", id))

	err = row.Scan(&task.Date, &task.Repeat)
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}
	if task.Date == "" {
		returnError(w, "задача не найдена", 500)
		return
	}

	// Удаление записи без правил повторения
	if task.Repeat == "" {
		del, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
			sql.Named("id", id),
		)
		if err != nil {
			returnError(w, err.Error(), 500)
			return
		}

		rows, err := del.RowsAffected()
		if err != nil {
			returnError(w, err.Error(), 500)
			return
		} else if rows == 0 {
			returnError(w, "задача не найдена", 500)
			return
		}
	} else {
		// Редактирование записи с правилами повторения
		task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			returnError(w, err.Error(), 500)
			return
		}

		upd, err := db.Exec("UPDATE scheduler SET date = :date WHERE id = :id",
			sql.Named("id", id),
			sql.Named("date", task.Date),
		)
		if err != nil {
			returnError(w, err.Error(), 500)
			return
		}

		rows, err := upd.RowsAffected()
		if err != nil {
			returnError(w, err.Error(), 500)
			return
		} else if rows == 0 {
			returnError(w, "задача не найдена", 500)
			return
		}
	}

	// Если всё правильно, возвращает пустой ответ
	resp, err := json.Marshal(struct{}{})
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		returnError(w, "не указан идентификатор", 400)
		return
	}

	// Получение нужной записи
	db, err := sql.Open("sqlite", DatabaseDir)
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	defer db.Close()

	// Удаление записи без правил повторения
	del, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id),
	)
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	rows, err := del.RowsAffected()
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	} else if rows == 0 {
		returnError(w, "задача не найдена", 500)
		return
	}

	// Если всё правильно, возвращает пустой ответ
	resp, err := json.Marshal(struct{}{})
	if err != nil {
		returnError(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
