package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/chnmk/todo-list-final-task/services"
	"github.com/chnmk/todo-list-final-task/tests"
)

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ResponseValid struct {
	Id string `json:"id"`
}

type ResponseInvalid struct {
	Error string `json:"error"`
}

var DatabaseDir = tests.DBFile

func TaskRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		taskRequestPOST(w, r)
	default:
		http.Error(w, "некорректный метод, ожидался POST", http.StatusMethodNotAllowed)
		return
	}
}

func taskRequestPOST(w http.ResponseWriter, r *http.Request) {
	var responseValid ResponseValid
	var responseInvalid ResponseInvalid
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Считывает данные
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		responseInvalid.Error = err.Error()
		returnInvalidPOST(w, responseInvalid, 400)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		responseInvalid.Error = err.Error()
		returnInvalidPOST(w, responseInvalid, 400)
		return
	}

	// Проверка наличия title
	if task.Title == "" {
		responseInvalid.Error = "отсутствует название задачи"
		returnInvalidPOST(w, responseInvalid, 400)
		return
	}

	// Если дата пустая или не указана, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	// Проверяет формат даты 20060102
	dateParsed, err := time.Parse("20060102", task.Date)
	if err != nil {
		responseInvalid.Error = "некорректный формат времени"
		returnInvalidPOST(w, responseInvalid, 400)
		return
	}

	// Если дата меньше сегодняшнего числа...
	if dateParsed.Before(time.Now()) {
		if task.Repeat != "" {
			// При указанном правиле повторения нужно вычислить и записать в таблицу дату выполнения
			task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				responseInvalid.Error = err.Error()
				returnInvalidPOST(w, responseInvalid, 400)
				return
			}
		} else {
			// Если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число
			task.Date = time.Now().Format("20060102")
		}

	}

	// Если всё правильно, добавляет новую запись
	db, err := sql.Open("sqlite", DatabaseDir)
	if err != nil {
		responseInvalid.Error = err.Error()
		returnInvalidPOST(w, responseInvalid, 500)
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
		responseInvalid.Error = err.Error()
		returnInvalidPOST(w, responseInvalid, 500)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		responseInvalid.Error = err.Error()
		returnInvalidPOST(w, responseInvalid, 500)
		return
	}

	// Если всё правильно, возвращает ответ
	responseValid.Id = strconv.FormatInt(id, 10)
	resp, err := json.Marshal(responseValid)
	if err != nil {
		responseInvalid.Error = err.Error()
		returnInvalidPOST(w, responseInvalid, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func returnInvalidPOST(w http.ResponseWriter, msg ResponseInvalid, status int) {
	resp, err := json.Marshal(msg)
	if err != nil {
		msg.Error = "ошибка при записи ответа"
		returnInvalidPOST(w, msg, 500)
	}

	if status == 400 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(resp)
}
