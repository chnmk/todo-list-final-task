package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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

type TasksArray struct {
	Tasks []Task `json:"tasks"`
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
		taskPOST(w, r)
	case http.MethodGet:
		taskGET(w, r)
	case http.MethodPut:
		taskPUT(w, r)
	default:
		returnError(w, "неожиданный метод запроса", 500)
		return
	}
}

func returnError(w http.ResponseWriter, msg string, status int) {
	var e ResponseInvalid
	e.Error = msg

	resp, err := json.Marshal(e)
	if err != nil {
		e.Error = "ошибка при записи ответа"
		returnError(w, msg, 500)
	}

	if status == 400 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(resp)
}

func checkRequest(w http.ResponseWriter, r *http.Request) (Task, error) {
	now := time.Now().Format("20060102")
	nowParsed, err := time.Parse("20060102", now)
	if err != nil {
		returnError(w, err.Error(), 500)
		return Task{}, err
	}

	var task Task
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		returnError(w, err.Error(), 500)
		return Task{}, err
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		returnError(w, err.Error(), 500)
		return Task{}, err
	}

	// Проверка наличия title
	if task.Title == "" {
		returnError(w, "отсутствует название задачи", 400)
		return Task{}, errors.New("отсутствует название задачи")
	}

	// Если дата пустая или не указана, берётся сегодняшнее число
	if task.Date == "" {
		task.Date = now
	}

	// Проверяет формат даты 20060102
	dateParsed, err := time.Parse("20060102", task.Date)
	if err != nil {
		returnError(w, "некорректный формат времени", 400)
		return Task{}, err
	}

	// Если дата меньше сегодняшнего числа...
	if dateParsed.Before(nowParsed) {
		if task.Repeat != "" {
			// При указанном правиле повторения нужно вычислить и записать в таблицу дату выполнения
			task.Date, err = services.NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				returnError(w, err.Error(), 500)
				return Task{}, err
			}
		} else {
			// Если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число
			task.Date = time.Now().Format("20060102")
		}

	}

	return task, nil
}