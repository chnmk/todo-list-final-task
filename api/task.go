package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/chnmk/todo-list-final-task/services"
	"github.com/chnmk/todo-list-final-task/tests"
)

type Task struct {
	Id      string `json:"id"`
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
	var responseInvalid ResponseInvalid

	switch r.Method {
	case http.MethodPost:
		taskPOST(w, r)
	case http.MethodGet:
		taskGET(w, r)
	case http.MethodPut:
		taskPUT(w, r)
	default:
		responseInvalid.Error = "неожиданный метод запроса"
		returnInvalid(w, responseInvalid, 500)
		return
	}
}

func taskPOST(w http.ResponseWriter, r *http.Request) {
	var responseValid ResponseValid
	var responseInvalid ResponseInvalid
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	now := time.Now().Format("20060102")
	nowParsed, err := time.Parse("20060102", now)
	if err != nil {
		responseInvalid.Error = err.Error()
		returnInvalid(w, responseInvalid, 500)
		return
	}

	// Считывает данные
	var task Task
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		responseInvalid.Error = err.Error()
		returnInvalid(w, responseInvalid, 400)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		responseInvalid.Error = err.Error()
		returnInvalid(w, responseInvalid, 400)
		return
	}

	// Проверка наличия title
	if task.Title == "" {
		responseInvalid.Error = "отсутствует название задачи"
		returnInvalid(w, responseInvalid, 400)
		return
	}

	// Если дата пустая или не указана, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = now
	}

	// Проверяет формат даты 20060102
	dateParsed, err := time.Parse("20060102", task.Date)
	if err != nil {
		responseInvalid.Error = "некорректный формат времени"
		returnInvalid(w, responseInvalid, 400)
		return
	}

	// Если дата меньше сегодняшнего числа...
	if dateParsed.Before(nowParsed) {
		if task.Repeat != "" {
			// При указанном правиле повторения нужно вычислить и записать в таблицу дату выполнения
			task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				responseInvalid.Error = err.Error()
				returnInvalid(w, responseInvalid, 400)
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
		returnInvalid(w, responseInvalid, 500)
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
		returnInvalid(w, responseInvalid, 500)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		responseInvalid.Error = err.Error()
		returnInvalid(w, responseInvalid, 500)
		return
	}

	// Если всё правильно, возвращает ответ
	responseValid.Id = strconv.FormatInt(id, 10)
	resp, err := json.Marshal(responseValid)
	if err != nil {
		responseInvalid.Error = err.Error()
		returnInvalid(w, responseInvalid, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func taskGET(w http.ResponseWriter, r *http.Request) {
	var responseInvalid ResponseInvalid
	var task Task

	id := r.FormValue("id")

	if id == "" {
		responseInvalid.Error = "не указан идентификатор"
		returnInvalid(w, responseInvalid, 500)
		return
	}

	// Подключение к базе
	db, err := sql.Open("sqlite", DatabaseDir)
	if err != nil {
		fmt.Println(4)
		responseInvalid.Error = err.Error()
		returnInvalid(w, responseInvalid, 500)
		return
	}

	defer db.Close()

	// Выполнение запроса
	rows, err := db.Query("SELECT * FROM scheduler WHERE id = :id LIMIT 1", sql.Named("id", id))
	if err != nil {
		responseInvalid.Error = err.Error()
		returnInvalid(w, responseInvalid, 500)
		return
	}

	defer rows.Close()

	// Чтение полученных данных
	for rows.Next() {
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			responseInvalid.Error = err.Error()
			returnInvalid(w, responseInvalid, 500)
			return
		}
	}

	// Запись ответа
	resp, err := json.Marshal(task)
	if err != nil {
		responseInvalid.Error = err.Error()
		returnInvalid(w, responseInvalid, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func taskPUT(w http.ResponseWriter, r *http.Request) {
	fmt.Println("placeholder")
}

func returnInvalid(w http.ResponseWriter, msg ResponseInvalid, status int) {
	resp, err := json.Marshal(msg)
	if err != nil {
		msg.Error = "ошибка при записи ответа"
		returnInvalid(w, msg, 500)
	}

	if status == 400 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(resp)
}
